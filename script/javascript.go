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
	vm        *otto.Otto
	sdk       string
	script    *bytes.Buffer
	processor Processor
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
	se.script.Reset()
	se.script.WriteString(se.sdk)
	se.script.WriteString(src)
}

func (se *JavaScriptImpl) AddScript(src string) {
	se.script.WriteString(src)
}

func (se *JavaScriptImpl) Run() (err error) {
	_, err = se.vm.Run(se.script.String())
	return
}

func New() *JavaScriptImpl {
	se := &JavaScriptImpl{
		vm:     otto.New(),
		sdk:    string(sdkSource),
		script: bytes.NewBufferString(""),
	}
	se.AddFunc("require", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return se
}

func NewWithProcessor(processor Processor) *JavaScriptImpl {
	se := New()
	se.processor = processor
	se.AddFunc("_http_get", HttpGet)
	se.AddFunc("_http_post", HttpPost)
	se.AddFunc("_http_delete", HttpDelete)
	se.AddFunc("_http_put", HttpPut)
	se.AddFunc("_sys_log", SysLog(se.processor))
	se.AddFunc("_sys_err", SysError(se.processor))
	se.AddFunc("require", SysRequire(se.processor))
	se.AddFunc("_resp_write", RespWrite(se.processor))
	se.AddFunc("_db_query", DBQuery(se.processor))
	se.AddFunc("_db_queryOne", DBQueryOne(se.processor))
	se.AddFunc("_db_save", DBSave(se.processor))
	se.AddFunc("_db_exec", DBExec(se.processor))
	return se
}

func init() {
	box := packr.NewBox("./")
	src, err := box.FindString(sdkFile)
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = []byte(src)
}
