package functions

import (
	"github.com/robertkrimen/otto"
	"github.com/zfd81/parrot/script"
)

func SysLog(env script.Environment) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			env.Println(arg.ToString())
		}
		return otto.Value{}
	}
}
