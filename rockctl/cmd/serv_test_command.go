package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	pb "github.com/zfd81/rock/proto/rockpb"

	"github.com/fatih/color"

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
		Errorf("open %s: No such file", path)
		return
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		Errorf(err.Error())
		return
	}

	request := &pb.RpcRequest{
		Header: map[string]string{"name": info.Name()},
		Data:   string(source),
	}
	resp, err := GetServiceClient().TestAnalysis(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}

	serv := &meta.Service{}
	err = json.Unmarshal([]byte(resp.Data), &serv)
	if err != nil {
		Errorf(err.Error())
		return
	}

	ps := map[string]string{}
	util.ReplaceBetween(serv.Path, "{", "}", func(i int, s int, e int, c string) (string, error) {
		name := strings.TrimSpace(c)
		v := readParameterInteractive(name)
		ps[name] = v
		return "", nil
	})

	for _, param := range serv.Params {
		v := readParameterInteractive(param.Name)
		ps[param.Name] = v
	}
	request.Params = ps

	resp, err = GetServiceClient().Test(context.Background(), request)
	if err != nil {
		Errorf(err.Error())
		return
	}
	fmt.Println("")
	color.Green("Service Definition:")
	fmt.Println("[NAMESPACE]", If(serv.Namespace == "", "default", serv.Namespace))
	fmt.Println("[PATH]", serv.Path)
	fmt.Println("[METHOD]", serv.Method)
	fmt.Println("")
	color.Green("Execution Results:")
	if len(resp.Header) > 0 {
		bytes, err := json.Marshal(resp.Header)
		if err != nil {
			Errorf(err.Error())
		}
		fmt.Println("[HEADER]", string(bytes))
	}
	fmt.Println("[DATA]", resp.Data)
	fmt.Println("")
	color.Green("Execution Logs:")
	fmt.Println(resp.Message)
	fmt.Println("Test complete.", "\n")
}

func readParameterInteractive(name string) string {
	var password string
	prompt := fmt.Sprintf("Value of %s: ", name)
	fmt.Print(prompt)
	fmt.Scanf("%s", &password)
	return password
}
