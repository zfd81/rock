package otto

import (
	"bytes"
	"fmt"
	"log"

	"github.com/zfd81/rock/core"

	"github.com/zfd81/rock/errs"

	"github.com/gobuffalo/packr/v2"

	js "github.com/robertkrimen/otto"
)

var (
	sdkFile   = "sdk.js"
	sdkSource string
)

type Function func(call js.FunctionCall) js.Value

type FuncResult struct {
	Normal     bool        `json:"-"`
	StatusCode int         `json:"code"`
	Data       interface{} `json:"data"`
	Message    string      `json:"msg"`
}

func Result(call js.FunctionCall, data interface{}) (value js.Value) {
	result := &FuncResult{
		Normal:     true,
		StatusCode: 200,
		Data:       data,
	}
	value, _ = call.Otto.ToValue(result)
	return
}

func ErrorResult(call js.FunctionCall, err string) (value js.Value) {
	result := &FuncResult{
		Normal:     false,
		StatusCode: 400,
		Message:    errs.ErrorStyleFunc(err),
	}
	value, _ = call.Otto.ToValue(result)
	return
}

type JavaScriptImpl struct {
	vm        *js.Otto
	sdk       string
	script    *bytes.Buffer
	processor core.Processor
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

func (se *JavaScriptImpl) AddFunc(name string, function interface{}) error {
	return se.vm.Set(name, function)
}

func (se *JavaScriptImpl) CallFunc(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
func (se *JavaScriptImpl) GetSdk() string {
	return se.sdk
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
		vm:     js.New(),
		sdk:    string(sdkSource),
		script: bytes.NewBufferString(""),
	}
	se.AddFunc("require", func(call js.FunctionCall) js.Value {
		return js.Value{}
	})
	return se
}

func NewWithProcessor(processor core.Processor) *JavaScriptImpl {
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
	se.AddFunc("_kv_get", KvGet(se.processor))
	se.AddFunc("_kv_set", KvSet(se.processor))
	return se
}

func GetSdk() string {
	return sdkSource
}

func throwException(format string, msgs ...interface{}) {
	err := fmt.Sprintf(format, msgs...)
	exception, _ := js.ToValue(err)
	panic(exception)
}

func init() {
	box := packr.New("sdk", "./")
	src, err := box.FindString(sdkFile)
	if err != nil {
		log.Fatal(err)
	}
	sdkSource = src
}
