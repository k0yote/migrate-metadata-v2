package k0yote3web

import "time"

const (
	publicIpfsGatewayUrl  = "https://ipfs.io/ipfs/"
	defaultIpfsGatewayUrl = "http://127.0.0.1:8080/ipfs/"

	defaultIpfsAPI = "http://127.0.0.1:5001"
	infuraAPI      = "https://ipfs.infura.io:5001"

	saveFolderName     = "internal"
	metadataFolderName = saveFolderName + "/" + "meta"
	imageFolderName    = saveFolderName + "/" + "image"
	uploadFolderName   = saveFolderName + "/" + "upload"

	fetchDownloadMetaLimit  = 300
	fetchDownloadImageLimit = 30
	waitTime                = 1 * time.Second
)
