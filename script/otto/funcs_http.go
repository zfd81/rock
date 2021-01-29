package otto

import (
	"log"
	"strings"

	"github.com/zfd81/rock/core"

	js "github.com/robertkrimen/otto"
	"github.com/zfd81/rock/httpclient"
)

func HttpGet(call js.FunctionCall) (value js.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	if strings.TrimSpace(url) == "" {
		throwException("Url cannot be empty")
	}
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			throwException(err.Error())
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
			throwException(err.Error())
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Get(url, data, ToHeader(header))
	value, _ = call.Otto.ToValue(*resp)
	return
}

func HttpPost(call js.FunctionCall) (value js.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	if strings.TrimSpace(url) == "" {
		throwException("Url cannot be empty")
	}
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			throwException(err.Error())
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
			throwException(err.Error())
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Post(url, data, ToHeader(header))
	value, _ = call.Otto.ToValue(*resp)
	return
}

func HttpDelete(call js.FunctionCall) (value js.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	if strings.TrimSpace(url) == "" {
		throwException("Url cannot be empty")
	}
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			throwException(err.Error())
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
			throwException(err.Error())
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Delete(url, data, ToHeader(header))
	value, _ = call.Otto.ToValue(*resp)
	return
}

func HttpPut(call js.FunctionCall) (value js.Value) {
	url := strings.TrimSpace(call.Argument(0).String())
	if strings.TrimSpace(url) == "" {
		throwException("Url cannot be empty")
	}
	var data map[string]interface{}
	var header map[string]interface{}
	data_v := call.Argument(1)
	if data_v.IsObject() {
		data_v, err := data_v.Export()
		if err != nil {
			throwException(err.Error())
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
			throwException(err.Error())
		} else {
			val, ok := header_v.(map[string]interface{})
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Put(url, data, ToHeader(header))
	value, _ = call.Otto.ToValue(*resp)
	return
}

func RespWrite(process core.Processor) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) js.Value {
		var data interface{}
		var err error

		data_v := call.Argument(0)
		if data_v.IsObject() {
			data, err = data_v.Export()
			if err != nil {
				log.Panicln(err)
				process.Println(err)
			}
		} else if data_v.IsString() {
			data, err = data_v.ToString()
			if err != nil {
				log.Panicln(err)
				process.Println(err)
			}
		} else if data_v.IsBoolean() {
			data, err = data_v.ToBoolean()
			if err != nil {
				log.Panicln(err)
				process.Println(err)
			}
		} else if data_v.IsNumber() {
			data, err = data_v.ToInteger()
			if err != nil {
				log.Panicln(err)
				process.Println(err)
			}
		}

		process.SetRespData(data)

		header_v := call.Argument(1)
		if header_v.IsObject() {
			header_v, err := header_v.Export()
			if err != nil {
				log.Panicln(err)
				process.Println(err)
			} else {
				val, ok := header_v.(map[string]interface{})
				if ok {
					for k, v := range val {
						process.AddRespHeader(k, v)
					}
				}
			}
		}
		return js.Value{}
	}
}

func ToHeader(h map[string]interface{}) httpclient.Header {
	header := httpclient.Header{}
	for k, v := range h {
		header.Set(k, v)
	}
	return header
}
