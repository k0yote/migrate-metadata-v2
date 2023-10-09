package k0yote3web

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	txPollTimeout             = 720
	txWaitTimeBetweenAttempts = time.Second * 1
	basefeeWiggleMultiplier   = 2
	txMaxAttempts             = 20
)

type contractHelper struct {
	*ProviderHandler
}

func newContractHelper(handler *ProviderHandler) (*contractHelper, error) {
	helper := &contractHelper{
		handler,
	}
	return helper, nil

}

func (c *contractHelper) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	zeroHash := common.Hash{}
	if txHash == zeroHash {
		return nil, fmt.Errorf("malformed transaction hash: [%v]", txHash)
	}

	var err error
	type ReceiptCh struct {
		ret *types.Receipt
		err error
	}

	provider := c.GetProvider()
	var timeoutFlag int32
	ch := make(chan *ReceiptCh, 1)

	go func() {
		for {
			receipt, err := provider.TransactionReceipt(ctx, txHash)
			if err != nil && err.Error() != "not found" {
				ch <- &ReceiptCh{
					err: err,
				}
				break
			}
			if receipt != nil {
				ch <- &ReceiptCh{
					ret: receipt,
					err: nil,
				}
				break
			}
			if atomic.LoadInt32(&timeoutFlag) == 1 {
				break
			}
		}
	}()

	select {
	case result := <-ch:
		if result.err != nil {
			return nil, err
		}

		return result.ret, nil
	case <-time.After(time.Duration(txPollTimeout) * txWaitTimeBetweenAttempts):
		atomic.StoreInt32(&timeoutFlag, 1)
		return nil, fmt.Errorf("transaction was not mined within %v seconds, "+
			"please make sure your transaction was properly sent. Be aware that it might still be mined", 720)
	}
}

func (c *contractHelper) awaitTx(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	provider := c.GetProvider()
	wait := txWaitTimeBetweenAttempts
	maxAttempts := uint8(txMaxAttempts)
	attempts := uint8(0)

	for {
		if attempts >= maxAttempts {
			return nil, fmt.Errorf("retry attempts to get tx exhausted, tx might have failed")
		}

		if tx, isPending, err := provider.TransactionByHash(ctx, hash); err != nil {
			attempts += 1
			time.Sleep(wait)
			continue
		} else {
			if isPending {
				time.Sleep(wait)
				continue
			}
			return tx, nil
		}
	}
}

func (c *contractHelper) DeployCreate2(factoryContract common.Address, contractABI, contractFactoryBin string, salt [32]byte, params ...interface{}) (common.Address, error) {

	parsed, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return common.Address{}, err
	}

	input, err := parsed.Pack("", params...)
	if err != nil {
		return common.Address{}, err
	}

	bytecode := common.FromHex(contractFactoryBin)

	addr := crypto.CreateAddress2(factoryContract, salt, crypto.Keccak256(append(bytecode, input...)))
	return addr, nil
}

// transact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
func (c *contractHelper) transact(opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	if opts.GasPrice != nil && (opts.GasFeeCap != nil || opts.GasTipCap != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// Create the transaction
	var (
		rawTx *types.Transaction
		err   error
	)
	if opts.GasPrice != nil {
		rawTx, err = c.createLegacyTx(opts, contract, input)
	} else if opts.GasFeeCap != nil && opts.GasTipCap != nil {
		rawTx, err = c.createDynamicTx(opts, contract, input, nil)
	} else {
		// Only query for basefee if gasPrice not specified
		if head, errHead := c.provider.HeaderByNumber(ensureContext(opts.Context), nil); errHead != nil {
			return nil, errHead
		} else if head.BaseFee != nil {
			rawTx, err = c.createDynamicTx(opts, contract, input, head)
		} else {
			// Chain is not London ready -> use legacy transaction
			rawTx, err = c.createLegacyTx(opts, contract, input)
		}
	}
	if err != nil {
		return nil, err
	}
	// Sign the transaction and schedule it for execution
	if opts.Signer == nil {
		return nil, fmt.Errorf("no signer to authorize the transaction with")
	}
	signedTx, err := opts.Signer(opts.From, rawTx)
	if err != nil {
		return nil, err
	}
	if opts.NoSend {
		return signedTx, nil
	}
	if err := c.sendTransaction(opts.Context, signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (c *contractHelper) sendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	return c.provider.SendTransaction(ensureContext(ctx), signedTx)
}

func (c *contractHelper) createLegacyTx(opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	if opts.GasFeeCap != nil || opts.GasTipCap != nil {
		return nil, errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
	}
	// Normalize value
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	// Estimate GasPrice
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		price, err := c.provider.SuggestGasPrice(ensureContext(opts.Context))
		if err != nil {
			return nil, err
		}
		gasPrice = price
	}
	// Estimate GasLimit
	gasLimit := opts.GasLimit
	if opts.GasLimit == 0 {
		var err error
		gasLimit, err = c.estimateGasLimit(opts, contract, input, gasPrice, nil, nil, value)
		if err != nil {
			return nil, err
		}
	}
	// create the transaction
	nonce, err := c.getNonce(opts)
	if err != nil {
		return nil, err
	}
	baseTx := &types.LegacyTx{
		To:       contract,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    value,
		Data:     input,
	}
	return types.NewTx(baseTx), nil
}

func (c *contractHelper) createDynamicTx(opts *bind.TransactOpts, contract *common.Address, input []byte, head *types.Header) (*types.Transaction, error) {
	// Normalize value
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	// Estimate TipCap
	gasTipCap := opts.GasTipCap
	if gasTipCap == nil {
		tip, err := c.provider.SuggestGasTipCap(ensureContext(opts.Context))
		if err != nil {
			return nil, err
		}
		gasTipCap = tip
	}
	// Estimate FeeCap
	gasFeeCap := opts.GasFeeCap
	if gasFeeCap == nil {
		gasFeeCap = new(big.Int).Add(
			gasTipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(basefeeWiggleMultiplier)),
		)
	}
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
	}
	// Estimate GasLimit
	gasLimit := opts.GasLimit
	if opts.GasLimit == 0 {
		var err error
		gasLimit, err = c.estimateGasLimit(opts, contract, input, nil, gasTipCap, gasFeeCap, value)
		if err != nil {
			return nil, err
		}
	}
	// create the transaction
	nonce, err := c.getNonce(opts)
	if err != nil {
		return nil, err
	}
	baseTx := &types.DynamicFeeTx{
		To:        contract,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		Value:     value,
		Data:      input,
	}
	return types.NewTx(baseTx), nil
}

func (c *contractHelper) estimateGasLimit(opts *bind.TransactOpts, contract *common.Address, input []byte, gasPrice, gasTipCap, gasFeeCap, value *big.Int) (uint64, error) {
	if contract != nil {
		addr := *contract
		// Gas estimation cannot succeed without code for method invocations.
		if code, err := c.provider.PendingCodeAt(ensureContext(opts.Context), addr); err != nil {
			return 0, err
		} else if len(code) == 0 {
			return 0, fmt.Errorf("gas estimation cannot succeed without code for method invocations")
		}
	}
	msg := ethereum.CallMsg{
		From:      opts.From,
		To:        contract,
		GasPrice:  gasPrice,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Value:     value,
		Data:      input,
	}
	return c.provider.EstimateGas(ensureContext(opts.Context), msg)
}

func (c *contractHelper) getNonce(opts *bind.TransactOpts) (uint64, error) {
	if opts.Nonce == nil {
		return c.provider.PendingNonceAt(ensureContext(opts.Context), opts.From)
	} else {
		return opts.Nonce.Uint64(), nil
	}
}

func (c *contractHelper) createRawTransaction(opts *bind.TransactOpts, to *common.Address) (*types.Transaction, error) {
	head, err := c.provider.HeaderByNumber(ensureContext(opts.Context), nil)
	if err != nil {
		return nil, err
	}

	gasTipCap, err := c.provider.SuggestGasTipCap(ensureContext(opts.Context))
	if err != nil {
		return nil, err
	}

	gasFeeCap := new(big.Int).Add(
		gasTipCap,
		new(big.Int).Mul(head.BaseFee, big.NewInt(basefeeWiggleMultiplier)),
	)

	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
	}

	estimatteFee := estimateGasFee(head.BaseFee, gasTipCap, gasFeeCap, Low)
	value := new(big.Int).Sub(opts.Value, estimatteFee)

	if value.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("not enough fund: gas price + value")
	}

	opts.Value = value

	nonce, err := c.getNonce(opts)
	if err != nil {
		return nil, err
	}
	baseTx := &types.DynamicFeeTx{
		To:        to,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       uint64(21_000),
		Value:     value,
		Data:      nil,
	}

	return types.NewTx(baseTx), nil
}

func estimateGasFee(baseFee, gasTipCap, gasFeeCap *big.Int, priority GasPriority) *big.Int {
	/*
		Here is an example:
		Your account balance: 1 ETH
		Estimated gas cost: 21,000 gas
		Current base fee: 20 Gwei
		Maximum priority fee you are willing to pay: 50 Gwei
		Maximum amount of Ether to send = 1 ETH - (21,000 gas * (20 Gwei + 50 Gwei)) = 1 ETH - 0.00147 ETH = 0.99853 ETH
	*/
	gasLimit := big.NewInt(21_000)
	return new(big.Int).Mul(gasLimit, new(big.Int).Add(baseFeeByPriority(baseFee, priority), gasTipCap))
}

func baseFeeByPriority(baseFee *big.Int, priority GasPriority) *big.Int {
	bf := FromWeiWithUnit(baseFee, EtherUnitGWei)
	x := new(big.Float).Mul(bf, big.NewFloat(priority.Value()))
	f, _ := x.Float64()
	return ToGWei(f)
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
