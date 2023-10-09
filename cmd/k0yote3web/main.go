package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	privateKey         string
	chainRpcUrl        string
	apiKey             string
	thirdpartyProvider string

	rootCmd = &cobra.Command{
		Use:   "k0yote3web",
		Short: "A CLI for the k0yote3web go SDK",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&privateKey, "privateKey", "k", "", "private key used to sign transactions")
	rootCmd.PersistentFlags().StringVarP(&chainRpcUrl, "chainRpcUrl", "u", "mumbai", "chain url where all rpc requests will be sent")
	rootCmd.PersistentFlags().StringVarP(&thirdpartyProvider, "thirdpartyProvider", "n", "alchemy", "third party provider")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apiKey", "a", "", "node provider api key")
	_ = viper.BindPFlag("privateKey", rootCmd.PersistentFlags().Lookup("privateKey"))
	_ = viper.BindPFlag("chainRpcUrl", rootCmd.PersistentFlags().Lookup("chainRpcUrl"))
	viper.SetDefault("chainRpcUrl", "polygon-mumbai")

	rootCmd.AddCommand(downloadCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := Execute(); err != nil {
		panic(err)
	}
}
