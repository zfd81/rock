package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"

	"github.com/zfd81/rock/meta"

	"github.com/spf13/cobra"
	pb "github.com/zfd81/rock/proto/rockpb"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <path> <method> [namespace]",
		Short: "Gets detailed information of a service",
		Run:   getCommandFunc,
	}
	return cmd
}

func getCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("get command requires service path and service method as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	request.Params["path"] = meta.FormatPath(args[0])
	request.Params["method"] = args[1]
	if len(args) > 2 {
		request.Params["namespace"] = cast.ToString(If(args[2] == "default", "", args[2]))
	}

	resp, err := GetServiceClient().FindService(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}

	var serv meta.Service
	err = json.Unmarshal([]byte(resp.Data), &serv)
	if err != nil {
		Errorf(err.Error())
		return
	}

	fmt.Println("+--------------+-----------+----------------------------------------------------+")
	fmt.Printf("%1s %12s %1s %9s %1s %50s %1s\n", "|", "NAMESPACE ", "|", "METHOD ", "|", "PATH                        ", "|")
	fmt.Println("+--------------+-----------+----------------------------------------------------+")
	if serv.Name != "" {
		fmt.Printf("%1s %12s %1s %9s %1s %50s %1s\n", "|", If(serv.Namespace == "", "default", serv.Namespace), "|", serv.Method, "|", serv.Path, "|")
	}
	fmt.Println("+--------------+-----------+----------------------------------------------------+")
}
