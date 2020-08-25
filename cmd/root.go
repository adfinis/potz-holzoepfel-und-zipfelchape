package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/pkg"
)

var (
	cfgFile           string
	listenAddr        string
	persistence       bool
	mongodbURI        string
	mongodbDatabase   string
	mongodbCollection string
	mongodbDocumentID string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "potz-holzoepfel-und-zipfelchape",
	Short: "Tri Tra Trulla La",
	Long:  "Dr Caasperli isch wider da! Dr Caasperli isch da.",
	Run: func(cmd *cobra.Command, args []string) {

		pkg.RunServer(listenAddr, persistence, mongodbURI, mongodbDatabase, mongodbCollection, mongodbDocumentID)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.potz-holzoepfel-und-zipfelchape.yaml)")
	rootCmd.PersistentFlags().StringVar(&listenAddr, "listen-addr", ":8080", "Listen address")
	rootCmd.PersistentFlags().BoolVar(&persistence, "persistence", false, "Enable persistence layer")
	rootCmd.PersistentFlags().StringVar(&mongodbURI, "mongodb-uri", "mongodb://root:hunter2@localhost:27017", "MongoDB URI")
	rootCmd.PersistentFlags().StringVar(&mongodbDatabase, "mongodb-database", "test", "MongoDB database")
	rootCmd.PersistentFlags().StringVar(&mongodbCollection, "mongodb-collection", "counter", "MongoDB collection")
	rootCmd.PersistentFlags().StringVar(&mongodbDocumentID, "mongodb-document-id", "DECAFBAD", "MongoDB counter document ID")
	if err := viper.BindPFlag("listen-addr", rootCmd.PersistentFlags().Lookup("listen-addr")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".potz-holzoepfel-und-zipfelchape" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".potz-holzoepfel-und-zipfelchape")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
