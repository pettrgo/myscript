package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"myscript/cmd/script"
	"myscript/config"
	"myscript/esmodel"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "my_script",
	Short: "简单脚本",
	Long:  "简单脚本",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string

func init() {
	rootCmd.AddCommand(script.OrderRefreshCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config/myscript.conf", "config file")
	fmt.Println(configFile)
	config.Init(configFile)
	cobra.OnInitialize(esmodel.InitEs)
}
