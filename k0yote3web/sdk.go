package k0yote3web

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type K0yote3WebSDK struct {
	*ProviderHandler
	Deployer *ContractDeployer
	Download Download
}

func NewK0yote3WebSDK(rpcUrlOrChainName string, options *SDKOptions) (*K0yote3WebSDK, error) {
	rpc, err := getDefaultRpcUrl(rpcUrlOrChainName, options.ApiKey, options.ThirdpartyProvier)
	if err != nil {
		return nil, err
	}

	provider, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	return NewThirdwebSDKFromProvider(provider, options)
}

func NewThirdwebSDKFromProvider(provider *ethclient.Client, options *SDKOptions) (*K0yote3WebSDK, error) {
	privateKey := ""

	// Override defaults with the options that are defined
	if options != nil {
		if options.PrivateKey != "" {
			privateKey = options.PrivateKey
		}
	}

	handler, err := NewProviderHandler(provider, privateKey)
	if err != nil {
		return nil, err
	}

	deployer, err := newContractDeployer(provider, privateKey)
	if err != nil {
		return nil, err
	}

	if deployer == nil {
		fmt.Println("Warning: Contract deployments are not supported on this network. SDK instantiated without a Deployer.")
	}

	sdk := &K0yote3WebSDK{
		ProviderHandler: handler,
		Deployer:        deployer,
	}

	return sdk, nil
}

// network to (e.g mainnet, sepolia, polygon, polygon-mumbai)
func getDefaultRpcUrl(network, apiKey string, thirdpartyProvider ThirdpartyProvider) (string, error) {
	switch thirdpartyProvider {
	case INFURA:
		return fmt.Sprintf("https://%s.infura.io/v3/%s", network, apiKey), nil
	case ALCHEMY:
		return fmt.Sprintf("https://%s.g.alchemy.com/v2/%s", network, apiKey), nil
	default:
		return "", fmt.Errorf("unknown third party provider: %s", thirdpartyProvider)
	}
}

func (sdk *K0yote3WebSDK) GetDownload(opts *DownloadMetaOptions) (*Download, error) {
	return newDownload(opts)
}

func (sdk *K0yote3WebSDK) GetRewriter(ipfsImageBaseURL, inputDir, outputDir string) (*MetaRewriter, error) {
	return newMetaRewriter(ipfsImageBaseURL, inputDir, outputDir)
}
