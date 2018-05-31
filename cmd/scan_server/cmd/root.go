/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"fmt"
	"os"
	"scan-api/database"
	"scan-api/log"
	"scan-api/server"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	serverConfigFile *string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server command ",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		var g errgroup.Group
		serverCfg, err := LoadConfigFromFile(*serverConfigFile)
		if err != nil {
			fmt.Printf("read config file failed %s", err.Error())
			return
		}

		if log.NewLogger(serverCfg.LogLevel, serverCfg.WriteLog) == nil {
			fmt.Println("Log init failed")
			return
		}

		if !database.InitDB() {
			fmt.Printf("init database error")
			return
		}

		scanServer := server.GetServer(&g, &serverCfg)
		if scanServer != nil {
			scanServer.RunServer()

			if serverCfg.SyncSwitch {
				database.BlockSync()
				database.StartSync(serverCfg.Interval)
			}

			if err := g.Wait(); err != nil {
				log.Error(err)
			}
		}

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
	serverConfigFile = rootCmd.Flags().StringP("config", "c", "", "server config file (required)")
	rootCmd.MarkFlagRequired("config")
}
