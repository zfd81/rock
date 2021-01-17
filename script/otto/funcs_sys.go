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

func SysRequire(ctx core.Context) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		path := strings.TrimSpace(call.Argument(0).String()) //获取依赖路径
		module := ctx.GetModule(meta.FormatPath(path))       //获取模块
		if module == nil || reflect.ValueOf(module).IsNil() {
			throwException("Module path[%s] not found", path)
		}
		_, err := call.Otto.Run("var module = {};" + module.GetSource())
		if err != nil {
			throwException(err.Error())
		}
		v, err := call.Otto.Get("module")
		if err != nil {
			throwException(err.Error())
		}
		if !v.IsObject() {
			throwException("Module %s definition error", path)
		}
		v, err = v.Object().Get("exports")
		if err != nil {
			throwException(err.Error())
		}
		if !v.IsObject() {
			throwException("Module %s definition error", path)
		}
		return v
	}
}
