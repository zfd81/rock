package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zfd81/parrot/http"
	"github.com/zfd81/parrot/meta"
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
		log.Println(prompt)
		return
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	ds := &meta.DataSource{}
	err = yaml.Unmarshal(yamlFile, ds)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Post(url("ds"), "application/json;charset=UTF-8", ds, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Content), &data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("[INFO] code:", data["code"])
			fmt.Println("[INFO] message:", data["msg"])
		}
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
		fmt.Println(err)
	} else {
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Content), &data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("[INFO] code:", data["code"])
			fmt.Println("[INFO] message:", data["msg"])
		}
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
		fmt.Println(err)
	} else {
		json, err := FormatJSON(resp.Content)
		if err != nil {
			fmt.Println(resp.Content)
		} else {
			fmt.Printf("[INFO] DataSource %s details:\n", name)
			fmt.Println(json)
		}
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
		fmt.Println(err)
	} else {
		json, err := FormatJSON(resp.Content)
		if err != nil {
			fmt.Println(resp.Content)
		} else {
			fmt.Println("[INFO] DataSourceice list:")
			fmt.Println(json)
		}
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
		log.Println(prompt)
		return
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	ds := &meta.DataSource{}
	err = yaml.Unmarshal(yamlFile, ds)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Put(url("ds"), ds, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Content), &data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("[INFO] code:", data["code"])
			fmt.Println("[INFO] message:", data["msg"])
		}
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
