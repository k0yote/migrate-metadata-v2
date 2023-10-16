package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	startTokenID int
	endTokenID   int
	baseURL      string
)

var downloadCmd = &cobra.Command{
	Use:   "download [command]",
	Short: "Interact with the Metadata from cloudstorage or self-hosted image distribution's service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please input a command to run")
	},
}

var downloadMetasCmd = &cobra.Command{
	Use:   "meta",
	Short: "meta data with download interface",
	Run: func(cmd *cobra.Command, args []string) {
		download, err := getDownload()
		if err != nil {
			panic(err)
		}

		if err := download.DownloadAndSaveMetadata(); err != nil {
			panic(err)
		}

		log.Printf("metadata download completed url: [%s] startTokenID: [%d] endTokenID: [%d] count: [%d]\n", baseURL, startTokenID, endTokenID, download.GetDownloadMetaCount())
	},
}

var downloadImagesCmd = &cobra.Command{
	Use:   "image",
	Short: "download image from metadata json",
	Run: func(cmd *cobra.Command, args []string) {
		download, err := getDownload()
		if err != nil {
			panic(err)
		}

		if err := download.DownloadAndSaveImage(); err != nil {
			panic(err)
		}

		log.Printf("images download completed count: [%d]\n", download.GetDownloadImageCount())
	},
}

func init() {
	downloadCmd.PersistentFlags().IntVarP(&startTokenID, "sTokenId", "s", 0, "start from download token id")
	downloadCmd.PersistentFlags().IntVarP(&endTokenID, "eTokenId", "e", 0, "end to download token id")
	downloadCmd.PersistentFlags().StringVarP(&baseURL, "baseUrl", "b", "", "base URL to download")

	downloadCmd.AddCommand(downloadMetasCmd)
	downloadCmd.AddCommand(downloadImagesCmd)
}
