package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zfd81/rooster/util"

	"github.com/spf13/cobra"
	"github.com/zfd81/rock/meta"
)

var (
	userShowDetail bool
)

// NewTestCommand returns the cobra command for "test".
func NewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <file> <param> (<param> can also be given from stdin)",
		Short: "Puts the given key into the store",
		Long:  ``,
		Run:   testCommandFunc,
	}
	return cmd
}

// testCommandFunc executes the "test" command.
func testCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("serv add command requires serv file as its argument"))
	}
	path := args[0]
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		prompt := fmt.Sprintf("open %s: No such file", path)
		fmt.Println("[ERROR]", prompt)
		return
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}
	params := map[string]interface{}{
		"name":   info.Name(),
		"source": string(source),
	}
	resp, err := client.Post(url("test/analysis"), "application/json;charset=UTF-8", params, nil)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}
	if resp.StatusCode != 200 {
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Content), &data)
		if err != nil {
			fmt.Println("[ERROR]", err)
		} else {
			fmt.Println("[ERROR] code:", data["code"])
			fmt.Println("[ERROR] message:", data["msg"], "\n")
		}
		return
	}

	serv := &meta.Service{}
	err = json.Unmarshal([]byte(resp.Content), &serv)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}
	ps := map[string]interface{}{}
	util.ReplaceBetween(serv.Path, "{", "}", func(i int, s int, e int, c string) (string, error) {
		name := strings.TrimSpace(c)
		v := readParameterInteractive(name)
		ps[name] = v
		return "", nil
	})
	for _, param := range serv.Params {
		v := readParameterInteractive(param.Name)
		if strings.ToUpper(param.DataType) == meta.DataTypeMap {
			m := map[string]interface{}{}
			if err := json.Unmarshal([]byte(v), &m); err != nil {
				Printerr("Parameter[" + param.Name + "] data type error")
				return
			}
			ps[param.Name] = m
		} else {
			ps[param.Name] = v
		}
	}
	params["params"] = ps
	resp, err = client.Post(url("test"), "application/json;charset=UTF-8", params, nil)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(resp.Content), &data)
	if err != nil {
		fmt.Println("[ERROR]", err)
	} else {
		fmt.Println("")
		fmt.Println("Service Definition:")
		fmt.Println("[NAMESPACE]", serv.Namespace)
		fmt.Println("[PATH]", serv.Path)
		fmt.Println("[METHOD]", serv.Method)
		fmt.Println("Execution Results:")
		fmt.Print(data["log"])
		if resp.StatusCode == 200 {
			bytes, err := json.Marshal(data["header"])
			if err != nil {
				fmt.Println("[HEADER]", data["header"])
			} else {
				if string(bytes) != "{}" {
					fmt.Println("[HEADER]", string(bytes))
				}
			}
			bytes, err = json.Marshal(data["data"])
			if err != nil {
				fmt.Println("[DATA]", data["data"])
			} else {
				if string(bytes) != "null" {
					fmt.Println("[DATA]", string(bytes))
				}
			}
		}
	}
	fmt.Println("Test complete.", "\n")
}

func readParameterInteractive(name string) string {
	var password string
	prompt := fmt.Sprintf("Value of %s: ", name)
	fmt.Print(prompt)
	fmt.Scanf("%s", &password)
	return password
}
