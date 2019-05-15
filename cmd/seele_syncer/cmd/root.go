/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/seeleteam/scan-api/database"
	"github.com/seeleteam/scan-api/log"
	"github.com/seeleteam/scan-api/syncer"
	"github.com/spf13/cobra"
)

var (
	serverConfigFile *string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "syncer command ",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		var g sync.WaitGroup
		serverCfg, err := LoadConfigFromFile(*serverConfigFile)
		if err != nil {
			fmt.Printf("read config file failed %s", err.Error())
			return
		}

		if log.NewLogger(serverCfg.LogFile, serverCfg.LogLevel, serverCfg.WriteLog) == nil {
			fmt.Println("Log init failed")
			return
		}

		dbClient := database.NewDBClient(serverCfg.DataBase, serverCfg.ShardNumber)
		if dbClient == nil {
			fmt.Printf("init database error")
			return
		}
		for i:=1; i<=4; i++ {
			dbClient.InitTxCntByShardNumber(i)
		}
		if serverCfg.DataBase.DataBaseMode == "replset" {
			dbClient.SetPrimaryMode()
		}

		syncer := syncer.NewSyncer(dbClient, serverCfg.RpcURL, serverCfg.ShardNumber)
		if syncer == nil {
			fmt.Printf("can not connect to node")
			return
		}

		syncer.StartSync(serverCfg.SyncInterval)
		g.Add(1)
		g.Wait()

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
