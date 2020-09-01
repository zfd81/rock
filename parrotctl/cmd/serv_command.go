package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zfd81/parrot/meta"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// NewServCommand returns the cobra command for "serv".
func NewServCommand() *cobra.Command {
	ac := &cobra.Command{
		Use:   "serv <subcommand>",
		Short: "Serv related commands",
	}
	ac.AddCommand(newServAddCommand())
	ac.AddCommand(newServDeleteCommand())
	ac.AddCommand(newServChangeCommand())
	ac.AddCommand(newServGetCommand())
	ac.AddCommand(newServListCommand())
	return ac
}

func newServAddCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add <file or directory>",
		Short: "Adds a new serv",
		Run:   servAddCommandFunc,
	}
	return &cmd
}

func newServDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <service method> <service path>",
		Short: "Deletes a serv",
		Run:   servDeleteCommandFunc,
	}
}

func newServChangeCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "change <serv file> [options]",
		Short: "Changes a serv",
		Run:   servChangeCommandFunc,
	}
	return &cmd
}

func newServGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get <service method> <service path>",
		Short: "Gets detailed information of a serv",
		Run:   servGetCommandFunc,
	}
	return &cmd
}

func newServListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all servs",
		Run:   servListCommandFunc,
	}
}

// servAddCommandFunc executes the "serv add" command.
func servAddCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv add command requires serv file as its argument"))
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
	serv := &meta.Service{}
	err = yaml.Unmarshal(yamlFile, serv)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Post(url("serv"), "application/json;charset=UTF-8", serv, nil)
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

// servDeleteCommandFunc executes the "serv delete" command.
func servDeleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv delete command requires service method and service path as its argument"))
	}
	method := args[0]
	path := meta.FormatPath(args[1])
	resp, err := client.Delete(url("serv/method/"+method+path), nil, nil)
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

// servGetCommandFunc executes the "serv get" command.
func servGetCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv get command requires service method and service path as its argument"))
	}
	method := args[0]
	path := meta.FormatPath(args[1])
	resp, err := client.Get(url("serv/method/"+method+path), nil)
	if err != nil {
		fmt.Println(err)
	} else {
		content := resp.Content
		var out bytes.Buffer
		err = json.Indent(&out, []byte(content), "", "  ")
		if err != nil {
			fmt.Println(resp.Content)
		} else {
			fmt.Println(fmt.Sprintf("[INFO] Service %s details:", path))
			fmt.Println(out.String())
		}
	}
}

// servListCommandFunc executes the "serv list" command.
func servListCommandFunc(cmd *cobra.Command, args []string) {
	path := "serv/list"
	if len(args) > 0 {
		path = path + meta.FormatPath(args[0])
	}
	resp, err := client.Get(url(path), nil)
	if err != nil {
		fmt.Println(err)
	} else {
		var data []string
		err = json.Unmarshal([]byte(resp.Content), &data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("[INFO] Service list:")
			for _, v := range data {
				fmt.Println(v)
			}
		}
	}
}

func servChangeCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv change command requires serv file as its argument"))
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
	serv := &meta.Service{}
	err = yaml.Unmarshal(yamlFile, serv)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Put(url("serv"), serv, nil)
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
