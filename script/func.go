package script

import (
	"log"

	"github.com/zfd81/parrot/httpclient"

	"github.com/robertkrimen/otto"
)

func post(call otto.FunctionCall) (value otto.Value) {
	url := call.Argument(0).String()
	var data map[string]interface{}
	var header map[string]string
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
			val, ok := header_v.(map[string]string)
			if ok {
				header = val
			}
		}
	}
	resp := httpclient.Post(url, data, header)
	value, _ = call.Otto.ToValue(*resp)
	return
}
