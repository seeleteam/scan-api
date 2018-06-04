/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"fmt"
	"os"
	"scan-api/chart/address"
	"scan-api/chart/block"
	"scan-api/chart/blockdiffculty"
	"scan-api/chart/blocktime"
	"scan-api/chart/hashrate"
	"scan-api/chart/topminers"
	"scan-api/chart/tx_history"
	"sync"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server command ",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		_, err := LoadConfigFromFile(*configFile)
		if err != nil {
			fmt.Printf("read config file failed %s", err.Error())
			return
		}
		go address.Process(&wg)
		wg.Add(1)
		go block.Process(&wg)
		wg.Add(1)
		go blockdifficulty.Process(&wg)
		wg.Add(1)
		go blocktime.Process(&wg)
		wg.Add(1)
		go hashrate.Process(&wg)
		wg.Add(1)
		go txhistory.Process(&wg)
		wg.Add(1)
		go topminers.Process(&wg)
		wg.Add(1)

		wg.Wait()
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
	configFile = rootCmd.Flags().StringP("config", "c", "", "config file (required)")
	rootCmd.MarkFlagRequired("config")
}
