package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "metrigocli",
	Short: "metrigocli is a client tools for metrigo-enabled apps",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
