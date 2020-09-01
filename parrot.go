package main

import (
	"fmt"

	"github.com/zfd81/parrot/conf"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/parrot/parrotctl/cmd"
	"github.com/zfd81/parrot/server/http"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:        "parrot",
		Short:      "parrot server",
		SuggestFor: []string{"parrot"},
		Run:        startCommandFunc,
	}
	port int
)

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8143, "Port to run the http http server")
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	r := gin.Default()
	parrot := r.Group("/parrot")
	api := parrot.Group("/api")
	{
		api.POST("/test", http.Test)
		api.POST("/serv", http.CreateService)
		api.DELETE("/serv/method/:method/*path", http.DeleteService)
		api.PUT("/serv", http.ModifyService)
		api.GET("/serv/method/:method/*path", http.FindService)
		api.GET("/serv/list/*path", http.ListService)
	}
	serv := parrot.Group("/serv")
	{
		serv.GET("/*path", http.CallGetService)
		serv.POST("/*path", http.CallPostService)
		serv.PUT("/*path", http.CallPutService)
		serv.DELETE("/*path", http.CallDeleteService)
	}
	r.Run(fmt.Sprintf(":%d", port))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.ExitWithError(cmd.ExitError, err)
	}
}

func main() {
	fmt.Println(conf.GetConfig().Banner)
	Execute()
}
