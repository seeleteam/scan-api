/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/seeleteam/scan-api/chart"
	_ "github.com/seeleteam/scan-api/chart/address"
	_ "github.com/seeleteam/scan-api/chart/block"
	_ "github.com/seeleteam/scan-api/chart/blockdifficulty"
	_ "github.com/seeleteam/scan-api/chart/blocktime"
	_ "github.com/seeleteam/scan-api/chart/hashrate"
	_ "github.com/seeleteam/scan-api/chart/topminers"
	_ "github.com/seeleteam/scan-api/chart/txhistory"
	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server command ",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		serverCfg, err := LoadConfigFromFile(*configFile)
		if err != nil {
			fmt.Printf("read config file failed %s", err.Error())
			return
		}

		if log.NewLogger(serverCfg.LogFile, serverCfg.LogLevel, serverCfg.WriteLog) == nil {
			fmt.Println("Log init failed")
			return
		}

		chart.GChartDB = database.NewDBClient(serverCfg.DataBase, 1)
		if chart.GChartDB == nil {
			fmt.Printf("init database error")
			return
		}

		chart.ShardCount = serverCfg.ShardCount
		processFuncs := chart.GetProcessFuncs()

		//start the all chart processor
		for i := 0; i < len(processFuncs); i++ {
			go processFuncs[i](&wg)
			wg.Add(1)
		}

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
