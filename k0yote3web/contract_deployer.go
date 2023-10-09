package k0yote3web

import "github.com/ethereum/go-ethereum/ethclient"

type ContractDeployer struct {
	*ProviderHandler
	helper *contractHelper
}

func newContractDeployer(provider *ethclient.Client, privateKey string) (*ContractDeployer, error) {
	handler, err := NewProviderHandler(provider, privateKey)
	if err != nil {
		return nil, err
	}

	helper, err := newContractHelper(handler)
	if err != nil {
		return nil, err
	}

	contractDeployer := &ContractDeployer{
		handler,
		helper,
	}

	return contractDeployer, nil
}
