package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zfd81/rock/server/api"

	pb "github.com/zfd81/rock/proto/rockpb"

	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"github.com/zfd81/rock/cluster"
	"github.com/zfd81/rock/conf"
	"github.com/zfd81/rock/server"

	"github.com/spf13/cobra"
	"github.com/zfd81/rock/rockctl/cmd"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	})
	//log.SetLevel(log.InfoLevel)
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	conf.GetConfig().Port = port
	conf.GetConfig().APIServer.Port = apiPort

	//打印配置信息
	log.Info("Rock version: ", conf.GetConfig().Version)
	log.Info("Rock etcd endpoints: ", conf.GetConfig().Etcd.Endpoints)

	RockServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.GetConfig().Port),
		Handler:      server.Router(),
		ReadTimeout:  conf.GetConfig().ReadTimeout * time.Second,
		WriteTimeout: conf.GetConfig().WriteTimeout * time.Second,
	}

	//开启RPC服务
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetConfig().APIServer.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := grpc.NewServer()
	pb.RegisterServiceServer(serv, &api.Service{})
	pb.RegisterDataSourceServer(serv, &api.DataSource{})
	pb.RegisterClusterServer(serv, &api.ClusterServer{})

	g.Go(func() error {
		err := serv.Serve(lis)
		if err != nil {
			log.Fatalf("failed to api serve: %v", err)
		}
		return err
	})
	g.Go(func() error {
		err := RockServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	server.WatchMeta()                  // 监测元数据变化
	server.InitDbs()                    // 根据元数据初始化数据源
	server.InitResources()              // 根据元数据初始化资源
	cluster.Register(time.Now().Unix()) // 集群注册
	time.Sleep(time.Duration(1) * time.Second)
	log.Infof("API server started successfully, listening on: %d", conf.GetConfig().APIServer.Port)
	log.Infof("Rock server started successfully, listening on: %d", conf.GetConfig().Port)

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
