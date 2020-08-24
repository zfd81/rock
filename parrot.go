package main

import (
	"fmt"

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
	}
	r.Run(fmt.Sprintf(":%d", port))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.ExitWithError(cmd.ExitError, err)
	}
}

func main() {
	Execute()
}
