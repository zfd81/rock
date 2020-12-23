package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	pb "github.com/zfd81/rock/proto/rockpb"

	"github.com/spf13/cast"

	"github.com/spf13/cobra"
	"github.com/zfd81/rock/meta"
)

// NewDataSourceCommand returns the cobra command for "ds".
func NewDataSourceCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "ds <subcommand>",
		Short: "DataSource related commands",
	}
	ac.AddCommand(newDataSourceAddCommand())
	ac.AddCommand(newDataSourceDeleteCommand())
	ac.AddCommand(newDataSourceGetCommand())
	ac.AddCommand(newDataSourceListCommand())
	return ac
}

func newDataSourceAddCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add <file or directory>",
		Short: "Adds a new datasource",
		Run:   dsAddCommandFunc,
	}
	return &cmd
}

func newDataSourceDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "del <name> [namespace]",
		Short: "Deletes a datasource",
		Run:   dsDeleteCommandFunc,
	}
}

func newDataSourceGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get <name> [namespace]",
		Short: "Gets detailed information of a datasource",
		Run:   dsGetCommandFunc,
	}
	return &cmd
}

func newDataSourceListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list [namespace]",
		Short: "Lists all datasources",
		Run:   dsListCommandFunc,
	}
}

// dsAddCommandFunc executes the "ds add" command.
func dsAddCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds add command requires datasource file as its argument"))
	}
	path := args[0]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		prompt := fmt.Sprintf("open %s: No such file", path)
		Printerr(prompt)
		return
	}
	definition, err := ioutil.ReadFile(path)
	if err != nil {
		Printerr(err.Error())
		return
	}
	request := &pb.RpcRequest{}
	request.Data = string(definition)
	resp, err := GetDataSourceClient().CreateDataSource(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	if resp.Code == 200 {
		Print(resp.Message)
	} else {
		Errorf(resp.Message)
	}
}

// dsDeleteCommandFunc executes the "ds delete" command.
func dsDeleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds del command requires datasource name as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	request.Params["name"] = args[0]
	if len(args) > 1 {
		request.Params["namespace"] = cast.ToString(If(args[1] == "default", "", args[1]))
	}
	resp, err := GetDataSourceClient().DeleteDataSource(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	if resp.Code == 200 {
		Print(resp.Message)
	} else {
		Errorf(resp.Message)
	}
}

// dsGetCommandFunc executes the "ds get" command.
func dsGetCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds get command requires datasource name as its argument"))
	}
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	request.Params["name"] = args[0]
	if len(args) > 1 {
		request.Params["namespace"] = cast.ToString(If(args[1] == "default", "", args[1]))
	}
	resp, err := GetDataSourceClient().FindDataSource(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var ds meta.DataSource
	err = json.Unmarshal([]byte(resp.Data), &ds)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("+--------------+-----------------+--------------+--------------+-----------+------------+------------+")
	fmt.Printf("%1s %12s %1s %15s %1s %12s %1s %12s %1s %9s %1s %10s %1s %10s %1s\n", "|", "NAMESPACE ", "|", "NAME     ", "|", "DRIVER   ", "|", "HOST    ", "|", "PORT  ", "|", "USER   ", "|", "DATABASE ", "|")
	fmt.Println("+--------------+-----------------+--------------+--------------+-----------+------------+------------+")
	if ds.Name != "" {
		fmt.Printf("%1s %12s %1s %15s %1s %12s %1s %12s %1s %9d %1s %10s %1s %10s %1s\n", "|", If(ds.Namespace == "", "default", ds.Namespace), "|", ds.Name, "|", ds.Driver, "|", ds.Host, "|", ds.Port, "|", ds.User, "|", ds.Database, "|")
	}
	fmt.Println("+--------------+-----------------+--------------+--------------+-----------+------------+------------+")
}

// dsListCommandFunc executes the "ds list" command.
func dsListCommandFunc(cmd *cobra.Command, args []string) {
	request := &pb.RpcRequest{}
	request.Params = map[string]string{}
	if len(args) > 0 {
		request.Params["namespace"] = cast.ToString(If(args[0] == "default", "", args[0]))
	}
	resp, err := GetDataSourceClient().ListDataSources(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	var dses []meta.DataSource
	err = json.Unmarshal([]byte(resp.Data), &dses)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("+-------+--------------+-----------------+--------------+--------------+-----------+------------+------------+")
	fmt.Printf("%1s %5s %1s %12s %1s %15s %1s %12s %1s %12s %1s %9s %1s %10s %1s %10s %1s\n", "|", "SEQ ", "|", "NAMESPACE ", "|", "NAME     ", "|", "DRIVER   ", "|", "HOST    ", "|", "PORT  ", "|", "USER   ", "|", "DATABASE ", "|")
	fmt.Println("+-------+--------------+-----------------+--------------+--------------+-----------+------------+------------+")
	for i, ds := range dses {
		fmt.Printf("%1s %5d %1s %12s %1s %15s %1s %12s %1s %12s %1s %9d %1s %10s %1s %10s %1s\n", "|", i+1, "|", If(ds.Namespace == "", "default", ds.Namespace), "|", ds.Name, "|", ds.Driver, "|", ds.Host, "|", ds.Port, "|", ds.User, "|", ds.Database, "|")
	}
	fmt.Println("+-------+--------------+-----------------+--------------+--------------+-----------+------------+------------+")

}

func FormatJSON(str string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(str), "", "  ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
