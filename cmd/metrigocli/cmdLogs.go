package main

import (
	"fmt"
	"strings"

	"github.com/mmelnyk/cvt"
	"github.com/mmelnyk/metrigocli/internal/metrigo"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Logs information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var logsshowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show logs levels settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := metrigo.NewClient(flagsHost)
		levels, latency, err := cli.GetLogLevels()

		fmt.Println(" Log levels for", flagsHost, ": ")
		fmt.Println(" Latency:"+cvt.BrWhiteFg, latency, cvt.ResetColor)

		if err != nil {
			fmt.Println(cvt.BrRedFg+"Failed"+cvt.ResetColor, err)
			return
		}

		for logger, level := range levels {
			fmt.Println(cvt.Tab, logger, ":", cvt.Tab, level)
		}

	},
}

var logssetCmd = &cobra.Command{
	Use:   "set <logger> <level>",
	Short: "Set logger level",
	Args:  cobra.ExactArgs(2),
	Long:  `Possible <level> values: FATAL, ERROR, WARNING, INFO, VERBOSE`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := metrigo.NewClient(flagsHost)

		levels := []string{"FATAL", "ERROR", "WARNING", "INFO", "VERBOSE"}
		// Check level name
		correctlevelart := false
		args[1] = strings.ToUpper(args[1])
		for _, v := range levels {
			if args[1] == v {
				correctlevelart = true
			}
		}

		if !correctlevelart {
			fmt.Println(cvt.BrRedFg+"Level argument ", args[1], " does not match any valid value"+cvt.ResetColor)
			return
		}

		fmt.Println(" Set logger", args[0], "to level", args[1], "for", flagsHost, ": ")
		latency, err := cli.SetLogLevel(args[0], args[1])

		fmt.Println(" Latency:"+cvt.BrWhiteFg, latency, cvt.ResetColor)

		if err != nil {
			fmt.Println(cvt.BrRedFg+"Failed"+cvt.ResetColor, err)
			return
		}
	},
}

func init() {
	logsCmd.AddCommand(logsshowCmd)
	logsCmd.AddCommand(logssetCmd)
	RootCmd.AddCommand(logsCmd)
	logsshowCmd.Flags().StringVarP(&flagsHost, "host", "H", "localhost:9110", "Metrigo host")
	logssetCmd.Flags().StringVarP(&flagsHost, "host", "H", "localhost:9110", "Metrigo host")
}
