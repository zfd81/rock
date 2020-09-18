package functions

import (
	"github.com/robertkrimen/otto"
	"github.com/zfd81/sunflower/script"
)

func SysLog(process script.Process) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			process.Println(arg.ToString())
		}
		return otto.Value{}
	}
}

func SysError(process script.Process) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			process.Perror(arg.ToString())
		}
		return otto.Value{}
	}
}
