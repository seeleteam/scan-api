/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package server

import (
	"net/http"

	"scan-api/log"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"scan-api/api/routers"
)

//ScanServer hold instance of server
type ScanServer struct {
	Server *http.Server
	G      *errgroup.Group
	config *Config
}

//initGin init gin engine
func initGin(config *Config) *gin.Engine {

	gin.SetMode(config.GinMode)

	if config.DisableConsoleColor {
		gin.DisableConsoleColor()
	}

	e := gin.New()

	// use logs middleware
	e.Use(log.Logger(log.GetLogger()))
	e.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		//AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-fc-version", "x-fc-terminal", "*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	corsConfig.AllowAllOrigins = true
	e.Use(cors.New(corsConfig))

	// By default, http.ListenAndServe (which gin.Run wraps) will serve an unbounded number of requests.
	// Limiting the number of simultaneous connections can sometimes greatly speed things up under load
	if config.LimitConnections > 0 {
		e.Use(limit.MaxAllowed(config.LimitConnections))
	}

	routers.InitRouters(e)

	return e
}

//GetServer init scanServer and return server instance
func GetServer(g *errgroup.Group, config *Config) (server *ScanServer) {
	if config == nil {
		log.Error("[initGin] server config is nil")
	}

	ginHandler := initGin(config)

	return &ScanServer{
		Server: &http.Server{
			Addr:           config.Addr,
			Handler:        ginHandler,
			ReadTimeout:    config.ReadTimeout * time.Second,
			WriteTimeout:   config.WriteTimeout * time.Second,
			IdleTimeout:    config.IdleTimeout * time.Second,
			MaxHeaderBytes: 1 << config.MaxHeaderBytes,
		},
		G: g,
	}
}

//RunServer Run our server in a goroutine
func (sl *ScanServer) RunServer() {
	sl.G.Go(func() error {
		return http.ListenAndServe(sl.Server.Addr, sl.Server.Handler)
	})
}
