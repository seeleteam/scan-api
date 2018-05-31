/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"fmt"
	"os"
	"scan-api/log"
	"scan-api/node"
	"sync"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server command ",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfigFromFile(*configFile)
		if err != nil {
			fmt.Printf("read config file failed %s", err.Error())
			return
		}

		if log.NewLogger(config.LogLevel, config.WriteLog) == nil {
			fmt.Println("Log init failed")
			return
		}

		var wg sync.WaitGroup
		node.StartFindNodeService(&config)
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
