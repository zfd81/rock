package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zfd81/parrot/http"
	"github.com/zfd81/parrot/meta"
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
		Use:   "del <path> <method> [namespace]",
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
		Use:   "get <path> <method> [namespace]",
		Short: "Gets detailed information of a serv",
		Run:   servGetCommandFunc,
	}
	return &cmd
}

func newServListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list [path] [namespace]",
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
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	param := map[string]interface{}{
		"name":   info.Name(),
		"source": string(source),
	}
	resp, err := client.Post(url("serv"), "application/json;charset=UTF-8", param, nil)
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
	if len(args) < 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv del command requires service path and service method as its argument"))
	}
	path := args[0]
	method := args[1]
	header := http.Header{}
	if len(args) > 2 {
		header.Set("namespace", args[2])
	}
	resp, err := client.Delete(url(fmt.Sprintf("serv/method/%s%s", method, meta.FormatPath(path))), nil, header)
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
	if len(args) < 2 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv get command requires service path and service method as its argument"))
	}
	path := args[0]
	method := args[1]
	header := http.Header{}
	if len(args) > 2 {
		header.Set("namespace", args[2])
	}
	resp, err := client.Get(url(fmt.Sprintf("serv/method/%s%s", method, meta.FormatPath(path))), nil, header)
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
	var path string
	header := http.Header{}
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
		header.Set("namespace", args[1])
		break
	}
	resp, err := client.Get(url(fmt.Sprintf("serv/list%s", meta.FormatPath(path))), nil, header)
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
				fmt.Printf("[INFO] %s \n", v)
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
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	param := map[string]interface{}{
		"name":   info.Name(),
		"source": string(source),
	}
	resp, err := client.Put(url("serv"), param, nil)
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
