package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"myscript/cmd/script"
	"myscript/config"
	"myscript/loggers"
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
	rootCmd.AddCommand(script.TestCmd)
	rootCmd.AddCommand(script.DayOrderShopSearchCmd)
	rootCmd.AddCommand(script.PubTransferMsgCmd)
	rootCmd.AddCommand(script.TaskSearchCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	config.Init(configFile)
	loggers.Init()
}
