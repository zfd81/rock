package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zfd81/rock/errs"

	"github.com/zfd81/rock/server/api"

	"github.com/spf13/cobra"
)

const (
	Version        = "1.0"
	cliName        = "rockctl"
	cliDescription = "A simple command line client for parrot."

	defaultDialTimeout      = 2 * time.Second
	defaultCommandTimeOut   = 5 * time.Second
	defaultKeepAliveTime    = 2 * time.Second
	defaultKeepAliveTimeOut = 6 * time.Second
)

var (
	globalFlags = GlobalFlags{}
	rootCmd     = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"rockctl"},
	}
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"127.0.0.1:8143"}, "gRPC endpoints")
	rootCmd.PersistentFlags().StringVar(&globalFlags.User, "user", "", "username[:password] for authentication (prompt if password is not supplied)")
	rootCmd.PersistentFlags().StringVar(&globalFlags.Password, "password", "", "password for authentication (if this option is used, --user option shouldn't include password)")

	rootCmd.AddCommand(
		NewVersionCommand(),
		NewGetCommand(),
		NewListCommand(),
		NewTestCommand(),
		NewServCommand(),
		NewDataSourceCommand(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ExitWithError(ExitError, err)
	}
}

func wrapResponse(content string) (*api.ApiResponse, error) {
	resp := &api.ApiResponse{}
	if err := json.Unmarshal([]byte(content), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func Printerr(msg string) {
	fmt.Printf("[ERROR] %s \n", errs.ErrorStyleFunc(msg))
}

func Print(format string, msgs ...interface{}) {
	fmt.Printf("[INFO] %s \n", fmt.Sprintf(format, msgs...))
}

func Errorf(format string, msgs ...interface{}) {
	fmt.Printf("[ERROR] %s \n", errs.ErrorStyleFunc(fmt.Sprintf(format, msgs...)))
}
