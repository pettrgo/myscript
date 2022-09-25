package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "my_script",
	Short: "简单脚本",
	Long:  "简单脚本",
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("pre run my script")
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}