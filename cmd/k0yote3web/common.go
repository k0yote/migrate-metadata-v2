package main

import (
	"github.com/thirdtool-dev/go-sdk/k0yote3web"
)

var (
	k0yote3webSDK *k0yote3web.K0yote3WebSDK
)

func initSdk() {
	if sdk, err := k0yote3web.NewK0yote3WebSDK(
		chainRpcUrl,
		&k0yote3web.SDKOptions{
			PrivateKey:        privateKey,
			ThirdpartyProvier: k0yote3web.ThirdpartyProvider(thirdpartyProvider),
			ApiKey:            apiKey,
		},
	); err != nil {
		panic(err)
	} else {
		k0yote3webSDK = sdk
	}
}

func getDownload() (*k0yote3web.Download, error) {
	if k0yote3webSDK == nil {
		initSdk()
	}

	return k0yote3webSDK.GetDownload(
		&k0yote3web.DownloadMetaOptions{
			BaseURL:      baseURL,
			StartTokenID: startTokenID,
			EndTokenID:   endTokenID,
		},
	)
}

func getRewrite() (*k0yote3web.MetaRewriter, error) {
	if k0yote3webSDK == nil {
		initSdk()
	}

	return k0yote3webSDK.GetRewriter(
		ipfsImageBaseURL,
		inputDir,
		outputDir,
	)
}
