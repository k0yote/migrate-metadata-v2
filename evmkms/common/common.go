package common

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type KmsSignature struct {
	R, S *big.Int
}

func KmsToEVMSignature(pubKey ecdsa.PublicKey,
	kmsSig KmsSignature,
	digestedMsg common.Hash,
) ([]byte, error) {
	if kmsSig.S.Cmp(CurveOrderHalf) > 0 {
		kmsSig.S = new(big.Int).Sub(CurveOrder, kmsSig.S)
	}

	if !ecdsa.Verify(&pubKey, digestedMsg[:], kmsSig.R, kmsSig.S) {
		return nil, fmt.Errorf("failed to verify signature")
	}

	pubKeyBytes := secp256k1.S256().Marshal(pubKey.X, pubKey.Y)

	rsSig := append(pad(kmsSig.R.Bytes(), 32), pad(kmsSig.S.Bytes(), 32)...)

	v := uint64(0)

	sig := append(rsSig, byte(v))
	recoveredPubKey, err := crypto.Ecrecover(digestedMsg[:], sig)
	if err != nil {
		return nil, fmt.Errorf("failed to recover pubKey with v = 0: %v", err)
	}

	if !bytes.Equal(recoveredPubKey, pubKeyBytes) {
		v = 1
		sig = append(rsSig, byte(v))
		recoveredPubKey, err = crypto.Ecrecover(digestedMsg[:], sig)
		if err != nil {
			return nil, fmt.Errorf("failed to recover pubKey with v = 1: %v", err)
		}
		if !bytes.Equal(recoveredPubKey, pubKeyBytes) {
			return nil, fmt.Errorf("cannot convert signature")
		}
	}

	return sig, nil
}

func pad(input []byte, paddedLength int) []byte {
	input = bytes.TrimLeft(input, "\x00")
	for len(input) < paddedLength {
		zeroBuf := []byte{0}
		input = append(zeroBuf, input...)
	}
	return input
}
