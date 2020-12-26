package goja

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/zfd81/rock/core"

	js "github.com/dop251/goja"
)

type rockscript struct {
	vm        *js.Runtime
	sdk       string
	script    *bytes.Buffer
	processor core.Processor
}

func (r *rockscript) AddVar(name string, value interface{}) error {
	r.vm.Set(name, value)
	return nil
}

func (r *rockscript) GetVar(name string) (interface{}, error) {
	value := r.vm.Get(name)
	switch value.ExportType().Kind() {
	case reflect.String:
		return value.ToString(), nil
	case reflect.Map:
		return value.Export(), nil
	case reflect.Struct:
		return value.Export(), nil
	case reflect.Int64:
		return value.ToInteger(), nil
	case reflect.Float64:
		return value.ToFloat(), nil
	case reflect.Bool:
		return value.ToBoolean(), nil
	default:
		return nil, nil
	}
}

func (r *rockscript) CallFunc(name string, args ...interface{}) (interface{}, error) {
	function, ok := js.AssertFunction(r.vm.Get(name))
	if !ok {
		return nil, fmt.Errorf("Function %s not found", name)
	}
	var params []js.Value
	if len(args) > 0 {
		for _, v := range args {
			params = append(params, r.vm.ToValue(v))
		}
	}
	value, err := function(js.Undefined(), params...)
	if err != nil {
		return nil, err
	}
	switch value.ExportType().Kind() {
	case reflect.String:
		return value.ToString(), nil
	case reflect.Map:
		return value.Export(), nil
	case reflect.Struct:
		return value.Export(), nil
	case reflect.Int64:
		return value.ToInteger(), nil
	case reflect.Float64:
		return value.ToFloat(), nil
	case reflect.Bool:
		return value.ToBoolean(), nil
	default:
		return nil, nil
	}
}

func (r *rockscript) AddFunc(name string, function interface{}) error {
	r.vm.Set(name, function)
	return nil
}

func (r *rockscript) GetSdk() string {
	return r.sdk
}

func (r *rockscript) SetScript(src string) {
	r.script.Reset()
	r.script.WriteString(r.sdk)
	r.script.WriteString(src)
}

func (r *rockscript) AddScript(src string) {
	r.script.WriteString(src)
}

func (se *rockscript) Run() (err error) {
	_, err = se.vm.RunString(se.script.String())
	return
}

func NewRockScript() *rockscript {
	se := &rockscript{
		vm:     js.New(),
		sdk:    string(""),
		script: bytes.NewBufferString(""),
	}
	return se
}
