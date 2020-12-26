package otto

import (
	"reflect"
	"strings"

	"github.com/zfd81/rock/script"

	"github.com/zfd81/rock/meta"

	js "github.com/robertkrimen/otto"
)

func SysLog(process script.Processor) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) js.Value {
		for _, arg := range call.ArgumentList {
			process.Println(arg.ToString())
		}
		return js.Value{}
	}
}

func SysError(process script.Processor) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) js.Value {
		for _, arg := range call.ArgumentList {
			process.Perror(arg.ToString())
		}
		return js.Value{}
	}
}

func SysRequire(process script.Processor) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		path := strings.TrimSpace(call.Argument(0).String())  //获取依赖路径
		module := process.SelectModule(meta.FormatPath(path)) //获取模块
		if module == nil || reflect.ValueOf(module).IsNil() {
			return script.ErrorResult(call, "Module path["+path+"] not found")
		}
		call.Otto.Set("exports", map[string]interface{}{})
		_, err := call.Otto.Run(module.GetSource())
		if err != nil {
			return js.NullValue()
		}
		m_v, err := call.Otto.Get("exports")
		if err != nil {
			return js.NullValue()
		}
		return m_v
	}
}
