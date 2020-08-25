package script

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/zfd81/parrot/http"

	"github.com/robertkrimen/otto"
)

func SysLog(env Environment) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			env.Println(arg.ToString())
		}
		return otto.Value{}
	}
}

func RespWrite(env Environment) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var data interface{}
		var err error

		data_v := call.Argument(0)
		if data_v.IsObject() {
			data, err = data_v.Export()
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			}
		} else if data_v.IsString() {
			data, err = data_v.ToString()
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			}
		} else if data_v.IsBoolean() {
			data, err = data_v.ToBoolean()
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			}
		} else if data_v.IsNumber() {
			data, err = data_v.ToInteger()
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			}
		}

		if data != nil {
			jsonStr, err := json.Marshal(data)
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			}
			env.SetRespContent(string(jsonStr))
		}

		header_v := call.Argument(1)
		if header_v.IsObject() {
			header_v, err := header_v.Export()
			if err != nil {
				log.Panicln(err)
				env.Println(err)
			} else {
				val, ok := header_v.(map[string]interface{})
				if ok {
					for k, v := range val {
						env.AddRespHeader(k, v)
					}
				}
			}
		}
		return otto.Value{}
	}
}

func HttpGet(call otto.FunctionCall) (value otto.Value) {
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
	resp := http.Get(url, data, header)
	value, _ = call.Otto.ToValue(*resp)
	return
}

func HttpPost(call otto.FunctionCall) (value otto.Value) {
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
	resp := http.Post(url, data, header)
	value, _ = call.Otto.ToValue(*resp)
	return
}
