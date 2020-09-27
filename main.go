package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zfd81/rock/server/api"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/zfd81/rock/cluster"
	"github.com/zfd81/rock/conf"
	"github.com/zfd81/rock/server"

	"github.com/spf13/cobra"
	"github.com/zfd81/rock/rockctl/cmd"
	"golang.org/x/sync/errgroup"
)

var (
	g       errgroup.Group
	rootCmd = &cobra.Command{
		Use:        "parrot",
		Short:      "parrot server",
		SuggestFor: []string{"parrot"},
		Run:        startCommandFunc,
	}
	port    int
	apiPort int
)

func init() {
	rootCmd.Flags().IntVar(&port, "port", conf.GetConfig().Port, "Port to run the server")
	rootCmd.Flags().IntVar(&apiPort, "api-port", conf.GetConfig().APIServer.Port, "Port to run the api server")
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	conf.GetConfig().Port = port
	conf.GetConfig().APIServer.Port = apiPort
	ApiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.GetConfig().APIServer.Port),
		Handler:      api.Router(),
		ReadTimeout:  conf.GetConfig().APIServer.ReadTimeout * time.Second,
		WriteTimeout: conf.GetConfig().APIServer.WriteTimeout * time.Second,
	}
	ParrotServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.GetConfig().Port),
		Handler:      server.Router(),
		ReadTimeout:  conf.GetConfig().ReadTimeout * time.Second,
		WriteTimeout: conf.GetConfig().WriteTimeout * time.Second,
	}
	g.Go(func() error {
		err := ApiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})
	g.Go(func() error {
		err := ParrotServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	server.WatchMeta()                  // 监测元数据变化
	server.InitDbs()                    // 根据元数据初始化数据源
	server.InitResources()              // 根据元数据初始化资源
	cluster.Register(time.Now().Unix()) // 集群注册
	color.Green("[INFO] API server listening on: %d", conf.GetConfig().APIServer.Port)
	color.Green("[INFO] Rock server listening on: %d", conf.GetConfig().Port)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.ExitWithError(cmd.ExitError, err)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	color.Green(conf.GetConfig().Banner)
	Execute()
}
