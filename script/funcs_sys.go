package script

import (
	"reflect"
	"strings"

	"github.com/zfd81/rock/meta"

	"github.com/robertkrimen/otto"
)

func SysLog(process Processor) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			process.Println(arg.ToString())
		}
		return otto.Value{}
	}
}

func SysError(process Processor) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			process.Perror(arg.ToString())
		}
		return otto.Value{}
	}
}

func SysRequire(process Processor) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) (value otto.Value) {
		path := strings.TrimSpace(call.Argument(0).String())  //获取依赖路径
		module := process.SelectModule(meta.FormatPath(path)) //获取模块
		if module == nil || reflect.ValueOf(module).IsNil() {
			return ErrorResult(call, "Module path["+path+"] not found")
		}
		call.Otto.Set("exports", map[string]interface{}{})
		_, err := call.Otto.Run(module.GetSource())
		if err != nil {
			return otto.NullValue()
		}
		m_v, err := call.Otto.Get("exports")
		if err != nil {
			return otto.NullValue()
		}
		return m_v
	}
}
