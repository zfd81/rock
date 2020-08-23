package script

import (
	"bytes"

	"github.com/spf13/cast"

	"github.com/robertkrimen/otto"
)

type JavaScriptImpl struct {
	vm     *otto.Otto
	log    *bytes.Buffer
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

func (se *JavaScriptImpl) Run() (string, error) {
	_, err := se.vm.Run(se.buffer.String())
	log := se.log.String()
	se.log.Reset()
	return log, err
}

func (se *JavaScriptImpl) Println(args ...interface{}) error {
	for _, arg := range args {
		se.log.WriteString(cast.ToString(arg))
	}
	se.log.WriteString("\n")
	return nil
}
