package script

import (
	"log"
	"strings"

	"github.com/zfd81/parrot/httpclient"

	"github.com/robertkrimen/otto"
)

func Log(se ScriptEngine) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			se.Println(arg.ToString())
		}
		return otto.Value{}
	}
}

func Get(call otto.FunctionCall) (value otto.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			log.Panicln(err)
		} else {
			val, ok := data_v.(map[string]interface{})
			if ok {
				data = val
			}
		}
	}
	header_v := call.Argument(2)
	if header_v.IsObject() {
		header_v, err := header_v.Export()
		if err != nil {
			log.Panicln(err)
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Get(url, data, header)
	value, _ = call.Otto.ToValue(*resp)
	return
}

func Post(call otto.FunctionCall) (value otto.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			log.Panicln(err)
		} else {
			val, ok := data_v.(map[string]interface{})
			if ok {
				data = val
			}
		}
	}
	header_v := call.Argument(2)
	if header_v.IsObject() {
		header_v, err := header_v.Export()
		if err != nil {
			log.Panicln(err)
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Post(url, data, header)
	value, _ = call.Otto.ToValue(*resp)
	return
}
