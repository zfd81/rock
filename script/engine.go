package script

import (
	"bytes"
	"log"

	"github.com/gobuffalo/packr"
	"github.com/robertkrimen/otto"
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
	box := packr.NewBox("./")
	src, err := box.FindString(sdkFile)
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = []byte(src)
}

func New() *JavaScriptImpl {
	se := &JavaScriptImpl{
		vm:     otto.New(),
		buffer: bytes.NewBuffer(sdkSource),
	}
	se.AddFunc("_post", post)
	return se
}
