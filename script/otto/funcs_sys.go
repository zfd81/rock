package otto

import (
	"reflect"
	"strings"

	"github.com/zfd81/rock/core"

	"github.com/zfd81/rock/meta"

	js "github.com/robertkrimen/otto"
)

func SysLog(process core.Processor) func(msgs []interface{}) {
	return func(msgs []interface{}) {
		process.Println(msgs...)
	}
}

func SysError(process core.Processor) func(msgs []interface{}) {
	return func(msgs []interface{}) {
		process.Perror(msgs...)
	}
}

func SysRequire(process core.Context) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		path := strings.TrimSpace(call.Argument(0).String()) //获取依赖路径
		module := process.GetModule(meta.FormatPath(path))   //获取模块
		if module == nil || reflect.ValueOf(module).IsNil() {
			throwException("Module path[%s] not found", path)
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
