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

var (
	sdkFile   = "sdk.js"
	sdkSource []byte
)

func init() {
	content, err := ioutil.ReadFile(sdkFile)
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
