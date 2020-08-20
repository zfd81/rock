package script

import (
	"bytes"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
)

type Function func(call otto.FunctionCall) otto.Value

type ScriptEngine interface {
	AddVar(name string, value interface{}) error
	AddFunc(name string, function Function) error
	SetScript(src string)
	AddScript(src string)
	Run() error
}

type JavaScriptImpl struct {
	vm     *otto.Otto
	buffer *bytes.Buffer
}

func (se *JavaScriptImpl) AddVar(name string, value interface{}) error {
	return se.vm.Set(name, value)
}

func (se *JavaScriptImpl) AddFunc(name string, function Function) error {
	return se.vm.Set(name, function)
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

var (
	sdkSource []byte
)

func init() {
	content, err := ioutil.ReadFile("sdk.js")
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = content
}

func New() *JavaScriptImpl {
	return &JavaScriptImpl{
		vm:     otto.New(),
		buffer: bytes.NewBuffer(sdkSource),
	}
}
