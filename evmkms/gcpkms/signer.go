package gcpkms

import (
	"context"
	"crypto/ecdsa"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"hash/crc32"
	"log"
	"math/big"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	common2 "github.com/thirdtool-dev/go-sdk/evmkms/common"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GoogleKMSClient struct {
	kmsClient *kms.KeyManagementClient
	ctx       context.Context
	cfg       Config
	publicKey *ecdsa.PublicKey
	signer    types.Signer
}

func NewGoogleKMSClient(ctx context.Context, cfg Config, txSigner ...types.Signer) (*GoogleKMSClient, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid config")
	}
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	// FIXME: chainId must be changed by setting
	signer := types.NewLondonSigner(new(big.Int).SetUint64(cfg.ChainID))
	if len(txSigner) > 0 {
		signer = txSigner[0]
	}

	c := &GoogleKMSClient{kmsClient: client, ctx: ctx, cfg: cfg, signer: signer}

	pubKey, err := c.getPublicKey()
	if err != nil {
		return nil, err
	}
	c.publicKey = pubKey

	return c, nil
}

func (c GoogleKMSClient) GetAddress() common.Address {
	return crypto.PubkeyToAddress(*c.publicKey)
}

func (c GoogleKMSClient) GetPublicKey() (*ecdsa.PublicKey, error) {
	return c.publicKey, nil
}

func (c GoogleKMSClient) SignHash(digest common.Hash) ([]byte, error) {
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)

	}
	digestCRC32C := crc32c(digest[:])

	req := &kmspb.AsymmetricSignRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
			c.cfg.ProjectID, c.cfg.LocationID, c.cfg.Key.Keyring, c.cfg.Key.Name, c.cfg.Key.Version),
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest[:],
			},
		},
		DigestCrc32C: wrapperspb.Int64(int64(digestCRC32C)),
	}

	result, err := c.kmsClient.AsymmetricSign(c.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sign digest: %v", err)
	}

	if !result.VerifiedDigestCrc32C {
		return nil, fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return nil, fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	return c.parseKMSSignature(digest, result.Signature)
}

func (c GoogleKMSClient) GetDefaultEVMTransactor() *bind.TransactOpts {
	return &bind.TransactOpts{
		Context: c.ctx,
		From:    c.GetAddress(),
		Signer:  c.GetEVMSignerFn(),
	}
}

func (c GoogleKMSClient) GetEVMSignerFn() bind.SignerFn {
	return func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if addr != c.GetAddress() {
			return nil, bind.ErrNotAuthorized
		}

		sig, err := c.SignHash(c.signer.Hash(tx))
		if err != nil {
			return nil, fmt.Errorf("cannot sign transaction: %v", err)
		}

		ret, err := tx.WithSignature(c.signer, sig)
		if err != nil {
			return nil, err
		}

		if _, err = c.HasSignedTx(ret); err != nil {
			return nil, err
		}

		return ret, nil
	}
}

func (c GoogleKMSClient) HasSignedTx(tx *types.Transaction) (bool, error) {
	from, err := types.Sender(c.signer, tx)
	if err != nil {
		return false, fmt.Errorf("cannot get sender of the tx: %v", err)
	}

	if from != c.GetAddress() {
		return false, fmt.Errorf("expected signer: %v, got %v", c.GetAddress(), from)
	}

	return true, nil
}

func (c *GoogleKMSClient) WithSigner(signer types.Signer) {
	c.signer = signer
}

func (c *GoogleKMSClient) WithChainID(chainID *big.Int) {
	if c.cfg.ChainID != chainID.Uint64() {
		c.cfg.ChainID = chainID.Uint64()
		c.signer = types.NewLondonSigner(chainID)
	}
}

func (c GoogleKMSClient) getPublicKey() (*ecdsa.PublicKey, error) {
	req := &kmspb.GetPublicKeyRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
			c.cfg.ProjectID, c.cfg.LocationID, c.cfg.Key.Keyring, c.cfg.Key.Name, c.cfg.Key.Version),
	}
	pubKey, err := c.kmsClient.GetPublicKey(c.ctx, req)
	if err != nil {
		return nil, err
	}

	return parseKMSPublicKey(pubKey)
}

func (c GoogleKMSClient) parseKMSSignature(digestedMsg common.Hash,
	kmsSignature []byte,
) ([]byte, error) {
	var sig common2.KmsSignature
	_, err := asn1.Unmarshal(kmsSignature, &sig)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kms signature: %v", err)
	}

	return common2.KmsToEVMSignature(*c.publicKey, sig, digestedMsg)
}

func (c GoogleKMSClient) Describe() error {
	listKeyRingsReq := &kmspb.ListKeyRingsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.cfg.ProjectID, c.cfg.LocationID),
	}

	it := c.kmsClient.ListKeyRings(c.ctx, listKeyRingsReq)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list Key rings: %v", err)
		}

		log.Printf("Key ring: %s\n", resp.Name)
	}

	return nil
}

func parseKMSPublicKey(kmsPubKey *kmspb.PublicKey) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(kmsPubKey.Pem))
	if block == nil || block.Type != "PUBLIC KEY" || len(block.Bytes) < 64 {
		return nil, fmt.Errorf("cannot decode public Key %v", kmsPubKey.Pem)
	}

	pubKeyBytes := block.Bytes[len(block.Bytes)-64:]
	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])

	if !secp256k1.S256().IsOnCurve(x, y) {
		return nil, fmt.Errorf("invalid secp256k1 public Key %v", kmsPubKey.Pem)
	}
	pubKey := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     x,
		Y:     y,
	}

	return &pubKey, nil
}
