package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cast"

	"github.com/zfd81/rock/http"
	"github.com/zfd81/rock/meta"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// NewDataSourceCommand returns the cobra command for "ds".
func NewDataSourceCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "ds <subcommand>",
		Short: "DataSource related commands",
	}
	ac.AddCommand(newDataSourceAddCommand())
	ac.AddCommand(newDataSourceDeleteCommand())
	ac.AddCommand(newDataSourceChangeCommand())
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

func newDataSourceChangeCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "change <file>",
		Short: "Changes a datasource",
		Run:   dsChangeCommandFunc,
	}
	return &cmd
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
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		Printerr(err.Error())
		return
	}
	ds := &meta.DataSource{}
	err = yaml.Unmarshal(yamlFile, ds)
	if err != nil {
		Printerr(err.Error())
		return
	}
	resp, err := client.Post(url("ds"), "application/json;charset=UTF-8", ds, nil)
	if err != nil {
		Printerr(err.Error())
		return
	}
	response, err := wrapResponse(resp.Content)
	if err != nil {
		Printerr(err.Error())
		return
	}
	if response.StatusCode == 200 {
		Print(response.Message)
	} else {
		Printerr(response.Message)
	}
}

// dsDeleteCommandFunc executes the "ds delete" command.
func dsDeleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds del command requires datasource name as its argument"))
	}
	name := args[0]
	header := http.Header{}
	if len(args) > 1 {
		header.Set("namespace", args[1])
	}
	resp, err := client.Delete(url(fmt.Sprintf("ds/name/%s", name)), nil, header)
	if err != nil {
		Printerr(err.Error())
		return
	}
	response, err := wrapResponse(resp.Content)
	if err != nil {
		Printerr(err.Error())
		return
	}
	if response.StatusCode == 200 {
		Print(response.Message)
	} else {
		Printerr(response.Message)
	}
}

// dsGetCommandFunc executes the "ds get" command.
func dsGetCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds get command requires datasource name as its argument"))
	}
	name := args[0]
	header := http.Header{}
	if len(args) > 1 {
		header.Set("namespace", args[1])
	}
	resp, err := client.Get(url(fmt.Sprintf("ds/name/%s", name)), nil, header)
	if err != nil {
		Printerr(err.Error())
		return
	}
	response, err := wrapResponse(resp.Content)
	if err != nil {
		Printerr(err.Error())
		return
	}
	if response.StatusCode == 200 {
		data := response.Data
		if data != nil {
			fmt.Printf("DataSource[%s] details:\n", name)
			ds := data.(map[string]interface{})
			fmt.Printf("%12s %15s %15s %15s %8s %8s %10s\n", "Namespace", "Name", "Driver", "Host", "Port", "User", "Database")
			fmt.Printf("%12s %15s %15s %15s %8s %8s %10s\n", ds["Namespace"], ds["Name"], ds["Driver"], ds["Host"], cast.ToString(ds["Port"]), ds["User"], ds["Database"])
		}
	} else {
		Printerr(response.Message)
	}
}

// dsListCommandFunc executes the "ds list" command.
func dsListCommandFunc(cmd *cobra.Command, args []string) {
	header := http.Header{}
	if len(args) > 0 {
		header.Set("namespace", args[0])
	}
	resp, err := client.Get(url("ds/list"), nil, header)
	if err != nil {
		Printerr(err.Error())
		return
	}
	response, err := wrapResponse(resp.Content)
	if err != nil {
		Printerr(err.Error())
		return
	}
	if response.StatusCode == 200 {
		data := response.Data
		if data != nil {
			dses, ok := data.([]interface{})
			if ok {
				fmt.Println("DataSource list:")
				fmt.Printf("%2s %12s %15s %15s %15s %8s %8s %10s\n", "", "Namespace", "Name", "Driver", "Host", "Port", "User", "Database")
				for i, v := range dses {
					ds := v.(map[string]interface{})
					fmt.Printf("%2d %12s %15s %15s %15s %8s %8s %10s\n", i, ds["Namespace"], ds["Name"], ds["Driver"], ds["Host"], cast.ToString(ds["Port"]), ds["User"], ds["Database"])
				}
			}
		}
	} else {
		Printerr(response.Message)
	}
}

func dsChangeCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("ds change command requires datasource file as its argument"))
	}
	path := args[0]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		prompt := fmt.Sprintf("open %s: No such file", path)
		Printerr(prompt)
		return
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		Printerr(err.Error())
		return
	}
	ds := &meta.DataSource{}
	err = yaml.Unmarshal(yamlFile, ds)
	if err != nil {
		Printerr(err.Error())
		return
	}
	resp, err := client.Put(url("ds"), ds, nil)
	response, err := wrapResponse(resp.Content)
	if err != nil {
		Printerr(err.Error())
		return
	}
	if response.StatusCode == 200 {
		Print(response.Message)
	} else {
		Printerr(response.Message)
	}
}

func FormatJSON(str string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(str), "", "  ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
