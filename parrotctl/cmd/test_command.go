package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/zfd81/parrot/httpclient"

	"github.com/spf13/cobra"
	"github.com/zfd81/parrot/meta"
	"gopkg.in/yaml.v2"
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

	for _, param := range serv.Params {
		v := readParameterInteractive(param.Name)
		param.Value = v
	}

	client := httpclient.New()
	url := fmt.Sprintf("http://%s/parrot/api/test", globalFlags.Endpoints[0])
	resp, err := client.Post(url, "application/json;charset=UTF-8", serv, nil)
	if err != nil {
		fmt.Println(err)
	}

	data := map[string]string{}
	err = json.Unmarshal([]byte(resp.Content), &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data["log"])
}

func readParameterInteractive(name string) string {
	var password string
	prompt := fmt.Sprintf("Value of %s: ", name)
	fmt.Print(prompt)
	fmt.Scanf("%s", &password)
	return password
}
