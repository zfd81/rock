package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	pb "github.com/zfd81/rock/proto/rockpb"
)

func NewChangeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change <serv file> [options]",
		Short: "Changes a serv",
		Run:   changeCommandFunc,
	}
	return cmd
}

func changeCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("change command requires service file as its argument"))
	}

	path := args[0]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		Errorf("open %s: No such file", path)
		return
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		Printerr(err.Error())
		return
	}

	request := &pb.RpcRequest{
		Params: map[string]string{"name": info.Name()},
		Data:   string(source),
	}
	resp, err := GetServiceClient().ModifyService(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	Print(resp.Message)
}
