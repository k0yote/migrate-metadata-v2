package k0yote3web

import "time"

const (
	defaultIpfsGatewayUrl = "https://ipfs.io/ipfs/"

	saveFolderName     = "internal"
	metadataFolderName = saveFolderName + "/" + "meta"
	imageFolderName    = saveFolderName + "/" + "image"
	uploadFolderName   = saveFolderName + "/" + "upload"

	fetchDownloadMetaLimit  = 300
	fetchDownloadImageLimit = 30
	waitTime                = 1 * time.Second
)
