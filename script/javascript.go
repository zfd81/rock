package script

import (
	"bytes"

	"github.com/robertkrimen/otto"
)

type JavaScriptImpl struct {
	vm     *otto.Otto
	buffer *bytes.Buffer
}

func (se *JavaScriptImpl) AddVar(name string, value interface{}) error {
	return se.vm.Set(name, value)
}

func (se *JavaScriptImpl) AddFunc(name string, function Function) error {
	return se.vm.Set(name, func(call otto.FunctionCall) otto.Value {
		return function(call)
	})
}

func (se *JavaScriptImpl) SetScript(src string) {
	se.buffer.Reset()
	se.buffer.WriteString(src)
}

func (se *JavaScriptImpl) AddScript(src string) {
	se.buffer.WriteString(src)
}

func (se *JavaScriptImpl) Run() (err error) {
	_, err = se.vm.Run(se.buffer.String())
	return
}
