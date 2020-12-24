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

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [path] [namespace]",
		Short: "Lists all services",
		Run:   listCommandFunc,
	}
	return cmd
}

func listCommandFunc(cmd *cobra.Command, args []string) {
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}

	var path string
	size := len(args)
	switch size {
	case 0:
		path = "/"
		break
	case 1:
		path = args[0]
		break
	default:
		path = args[0]
		request.Params["namespace"] = cast.ToString(If(args[1] == "default", "", args[1]))
		break
	}
	request.Params["path"] = path
	resp, err := GetServiceClient().ListServices(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}

	var servs []meta.Service
	err = json.Unmarshal([]byte(resp.Data), &servs)
	if err != nil {
		Errorf(err.Error())
		return
	}

	fmt.Println("+-------+--------------+-----------+----------------------------------------------------+")
	fmt.Printf("%1s %5s %1s %12s %1s %9s %1s %50s %1s\n", "|", "SEQ ", "|", "NAMESPACE ", "|", "METHOD ", "|", "PATH                        ", "|")
	fmt.Println("+-------+--------------+-----------+----------------------------------------------------+")
	for i, serv := range servs {
		fmt.Printf("%1s %5d %1s %12s %1s %9s %1s %50s %1s\n", "|", i+1, "|", If(serv.Namespace == "", "default", serv.Namespace), "|", serv.Method, "|", serv.Path, "|")
	}
	fmt.Println("+-------+--------------+-----------+----------------------------------------------------+")
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
