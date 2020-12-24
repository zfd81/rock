package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cast"

	"github.com/zfd81/rock/meta"

	"github.com/spf13/cobra"
	pb "github.com/zfd81/rock/proto/rockpb"
)

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del <path> <method> [namespace]",
		Short: "Deletes a service",
		Run:   deleteCommandFunc,
	}
	return cmd
}

func deleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("del command requires service path and service method as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	request.Params["path"] = meta.FormatPath(args[0])
	request.Params["method"] = args[1]
	if len(args) > 2 {
		request.Params["namespace"] = cast.ToString(If(args[2] == "default", "", args[2]))
	}
	resp, err := GetServiceClient().DeleteService(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	Print(resp.Message)
}
