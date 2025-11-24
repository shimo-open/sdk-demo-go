package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Config file path
var config string

func init() {
	confEnv := os.Getenv("EGO_CONFIG_PATH")
	if confEnv == "" {
		confEnv = "config/local.toml"
	}
	RootCommand.PersistentFlags().StringVarP(&config, "config", "c", confEnv, "Specify the config file (default config/local.toml)")
}

var RootCommand = &cobra.Command{
	Use: "sdk-demo-go",
}
