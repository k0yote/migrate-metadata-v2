package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var (
	CurveOrder     = secp256k1.S256().Params().N
	CurveOrderHalf = new(big.Int).Div(CurveOrder, new(big.Int).SetUint64(2))
)
