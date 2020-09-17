package script

import (
	"bytes"
	"log"

	"github.com/gobuffalo/packr"

	"github.com/robertkrimen/otto"
)

var (
	sdkFile   = "sdk.js"
	sdkSource []byte
)

type JavaScriptImpl struct {
	vm     *otto.Otto
	sdk    string
	buffer *bytes.Buffer
}

func (se *JavaScriptImpl) AddVar(name string, value interface{}) error {
	return se.vm.Set(name, value)
}

func (se *JavaScriptImpl) GetVar(name string) (interface{}, error) {
	value, err := se.vm.Get(name)
	if err != nil {
		return nil, err
	}
	if value.IsString() {
		return value.ToString()
	} else if value.IsObject() {
		return value.Export()
	} else if value.IsNumber() {
		val, err := value.ToInteger()
		if err != nil {
			return value.ToFloat()
		}
		return val, nil
	} else if value.IsBoolean() {
		return value.ToBoolean()
	}
	return nil, nil
}

func (se *JavaScriptImpl) AddFunc(name string, function Function) error {
	return se.vm.Set(name, func(call otto.FunctionCall) otto.Value {
		return function(call)
	})
}

func (se *JavaScriptImpl) SetScript(src string) {
	se.buffer.Reset()
	se.buffer.WriteString(se.sdk)
	se.buffer.WriteString(src)
}

func (se *JavaScriptImpl) AddScript(src string) {
	se.buffer.WriteString(src)
}

func (se *JavaScriptImpl) Run() (err error) {
	_, err = se.vm.Run(se.buffer.String())
	return
}

func New() ScriptEngine {
	return &JavaScriptImpl{
		vm:     otto.New(),
		sdk:    string(sdkSource),
		buffer: bytes.NewBuffer(sdkSource),
	}
}

func init() {
	box := packr.NewBox("./")
	src, err := box.FindString(sdkFile)
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = []byte(src)
}
