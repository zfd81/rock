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

type Function func(call otto.FunctionCall) otto.Value

type ScriptEngine interface {
	AddVar(name string, value interface{}) error
	AddFunc(name string, function Function) error
	SetScript(src string)
	AddScript(src string)
	Run() error
}

type Environment interface {
	Println(args ...interface{}) error
	SetRespStatus(code int)
	AddRespHeader(name string, value interface{})
	SetRespData(data interface{})
}

func init() {
	box := packr.NewBox("./")
	src, err := box.FindString(sdkFile)
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = []byte(src)
}

func New() ScriptEngine {
	return &JavaScriptImpl{
		vm:     otto.New(),
		sdk:    string(sdkSource),
		buffer: bytes.NewBuffer(sdkSource),
	}
}
