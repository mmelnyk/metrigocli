package main

import (
	"archive/tar"
	"fmt"
	"os"
	"time"

	"github.com/mmelnyk/cvt"
	"github.com/mmelnyk/metrigocli/internal/metrigo"
	"github.com/spf13/cobra"
)

var profilerCmd = &cobra.Command{
	Use:   "profiler",
	Short: "Manage profiling options",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var pprofcollectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect pprof profiles",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cli := metrigo.NewClient(flagsHost)

		var collection = []struct {
			name, file string
			call       func() ([]byte, time.Duration, error)
		}{
			{"metrics A", "metrics@start.json", cli.GetBlobMetrics},
			{"heap A", "heap@start", cli.GetHeap},
			{"mutex", "mutex", cli.GetMutex},
			{"block", "block", cli.GetBlock},
			{"trace", "trace.out", cli.GetTrace},
			{"CPU profiling", "cpu", cli.GetProfile},
			{"goroutines", "goroutines", cli.GetGoroutine},
			{"threadcreate", "threadcreate", cli.GetThreadCreate},
			{"heap B", "heap@finish", cli.GetHeap},
			{"metrics B", "metrics@finish.json", cli.GetBlobMetrics},
		}

		var tr *tar.Writer

		if f, err := os.Create(flagsOutput); err != nil {
			fmt.Println(cvt.BrRedFg+"Failed"+cvt.ResetColor, err)
			return
		} else {
			tr = tar.NewWriter(f)
		}

		defer tr.Close()

		fmt.Println(" Collecting pprof profiles from", flagsHost, " to ", flagsOutput, "...")

		for _, pprof := range collection {
			fmt.Print("    Collecting  ", pprof.name, "...")
			if val, latency, err := pprof.call(); err == nil {
				fmt.Println(" Latency:"+cvt.BrWhiteFg, latency, cvt.ResetColor)

				hdr := &tar.Header{
					Name: pprof.file,
					Mode: 0600,
					Size: int64(len(val)),
				}
				if err := tr.WriteHeader(hdr); err != nil {
					fmt.Println("   ", cvt.BrRedFg+"Warning "+cvt.ResetColor, err)
				}

				if _, err := tr.Write(val); err != nil {
					fmt.Println("   ", cvt.BrRedFg+"Warning "+cvt.ResetColor, err)
				}
			} else {
				fmt.Println(cvt.BrRedFg+"Failed"+cvt.ResetColor, err)
			}
		}
		fmt.Println(" Done")
	},
}

func init() {
	profilerCmd.AddCommand(pprofcollectCmd)
	RootCmd.AddCommand(profilerCmd)
	pprofcollectCmd.Flags().StringVarP(&flagsHost, "host", "H", "localhost:9110", "Metrigo host")
	pprofcollectCmd.Flags().StringVarP(&flagsOutput, "output", "O", "pprof.tar", "Output file")
}
