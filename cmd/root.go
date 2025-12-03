package cmd

import (
	"log"
	"os"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/spf13/cobra"
)

// Config file path
var Config string

func init() {
	confEnv := os.Getenv("EGO_CONFIG_PATH")
	if confEnv == "" {
		confEnv = "config/local.toml"
	}
	RootCommand.PersistentFlags().StringVarP(&Config, "config", "c", confEnv, "Specify the config file (default config/local.toml)")
}

var RootCommand = &cobra.Command{
	Use: "sdk-demo-go",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Println("ConfigFile", Config)
		provider, parser, tag, err := manager.NewDataSource(Config, true)
		if err != nil {
			log.Fatal("load config fail, ", err)
		}
		if err := econf.LoadFromDataSource(provider, parser, econf.WithSquash(true), econf.WithTagName(tag)); err != nil {
			log.Fatal("data source: load config, unmarshal config err", err)
		}
	},
}
