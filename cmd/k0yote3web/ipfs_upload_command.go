package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	providerType string
	projectID    string
	secret       string
	filepath     string
)

var ipfsUploadCmd = &cobra.Command{
	Use:   "ipfs-upload [command]",
	Short: "Upload meta file or meta files in directory to ipfs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please input a command to run")
	},
}

var ipfsUploadMetasCmd = &cobra.Command{
	Use:   "meta",
	Short: "meta data with upload interface",
	Run: func(cmd *cobra.Command, args []string) {
		ipfsUpload, err := getIpfsUploader()
		if err != nil {
			panic(err)
		}

		cid, err := ipfsUpload.Upload(filepath)
		if err != nil {
			panic(err)
		}

		log.Printf("CID: [%s] gatewayUrl: [%s]\n", cid.String(), ipfsUpload.GetGatewayUrl()+cid.String())
	},
}

func init() {
	ipfsUploadCmd.PersistentFlags().StringVarP(&providerType, "providerType", "t", "local", "ipfs api provider type to (e.g. local or infura)")
	ipfsUploadCmd.PersistentFlags().StringVarP(&projectID, "projectId", "p", "", "api projectId for using infura")
	ipfsUploadCmd.PersistentFlags().StringVarP(&secret, "secret", "s", "", "api secret for using infura")
	ipfsUploadCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "upload file or directory path for upload")

	ipfsUploadCmd.AddCommand(ipfsUploadMetasCmd)
}
