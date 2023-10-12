package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	ipfsImageBaseURL, inputDir, outputDir string
)

var rewriteCmd = &cobra.Command{
	Use:   "rewrite [command]",
	Short: "Interact with the Metadata from cloudstorage or self-hosted service interface",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please input a command to run")
	},
}

var rewriteMetaCmd = &cobra.Command{
	Use:   "replace-meta",
	Short: "rewrite image url with ipfs in metadata json",
	Run: func(cmd *cobra.Command, args []string) {
		rewriter, err := getRewrite()
		if err != nil {
			panic(err)
		}

		if err := rewriter.Rewrite(); err != nil {
			panic(err)
		}

		log.Printf("replaced image urls with ipfs urls count: [%d]\n", rewriter.Counter())
	},
}

func init() {
	rewriteCmd.PersistentFlags().StringVarP(&ipfsImageBaseURL, "ipfsImageBaseUrl", "g", "", "ipfs image Base URL")
	rewriteCmd.PersistentFlags().StringVarP(&inputDir, "inputDir", "i", "", "the folder of metadata files")
	rewriteCmd.PersistentFlags().StringVarP(&outputDir, "outputDir", "o", "", "the output folder of replaced image urls with ipfs")

	rewriteCmd.AddCommand(rewriteMetaCmd)
}
