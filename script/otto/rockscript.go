package otto

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/zfd81/rock/core"

	"github.com/gobuffalo/packr/v2"

	"github.com/robertkrimen/otto"
	js "github.com/robertkrimen/otto"
)

var (
	sdkFile   = "sdk.js"
	sdkSource string
)

type RockFunction struct {
	name     string
	function otto.Value
}

func (f RockFunction) Name() string {
	return f.name
}

func (f RockFunction) Perform(args ...interface{}) (interface{}, error) {
	value, err := f.function.Call(f.function, args...)
	if err != nil {
		return nil, err
	}
	if value.IsString() {
		return value.ToString()
	}
	if value.IsObject() {
		return value.Export()
	}
	if value.IsNumber() {
		val, err := value.ToInteger()
		if err != nil {
			return value.ToFloat()
		}
		return val, nil
	}
	if value.IsBoolean() {
		return value.ToBoolean()
	}
	return nil, fmt.Errorf("Method %s cannot find the data type of the return value", f.name)
}

type JavaScriptImpl struct {
	vm      *js.Otto
	sdk     string
	script  *bytes.Buffer
	context core.Context
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

func (se *JavaScriptImpl) GetMlVar(name string) (interface{}, error) {
	names := strings.Split(name, ".")
	value, err := se.vm.Get(names[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(names); i++ {
		if !value.IsObject() {
			return nil, nil
		}
		value, err = value.Object().Get(names[i])
		if err != nil {
			return nil, err
		}
	}
	if value.IsString() {
		return value.ToString()
	}
	if value.IsObject() {
		return value.Export()
	}
	if value.IsNumber() {
		val, err := value.ToInteger()
		if err != nil {
			return value.ToFloat()
		}
		return val, nil
	}
	if value.IsBoolean() {
		return value.ToBoolean()
	}
	if value.IsUndefined() || value.IsNull() {
		return nil, nil
	}
	return nil, fmt.Errorf("The data type of variable[%s] was not found", name)
}

func (se *JavaScriptImpl) AddFunc(name string, function interface{}) error {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		return fmt.Errorf("Wrong argument type: %s is not a function", name)
	}
	return se.vm.Set(name, function)
}

func (se *JavaScriptImpl) GetFunc(name string) (core.Function, error) {
	v, err := se.vm.Get(name)
	if err != nil {
		return nil, err
	}
	if !v.IsFunction() {
		return nil, fmt.Errorf("%s is not a function", name)
	}
	return &RockFunction{
		name:     name,
		function: v,
	}, nil
}

func (se *JavaScriptImpl) GetMlFunc(name string) (core.Function, error) {
	names := strings.Split(name, ".")
	value, err := se.vm.Get(names[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(names); i++ {
		if !value.IsObject() {
			return nil, nil
		}
		value, err = value.Object().Get(names[i])
		if err != nil {
			return nil, err
		}
	}
	if !value.IsFunction() {
		return nil, nil
	}
	return &RockFunction{
		name:     name,
		function: value,
	}, nil
}

func (se *JavaScriptImpl) CallFunc(name string, args ...interface{}) (interface{}, error) {
	f, err := se.GetFunc(name)
	if err != nil {
		return nil, err
	}
	return f.Perform(args...)
}

func (se *JavaScriptImpl) GetSdk() string {
	return se.sdk
}

func (se *JavaScriptImpl) SetScript(src string) {
	se.script.Reset()
	se.script.WriteString(src)
}

func (se *JavaScriptImpl) AddScript(src string) {
	se.script.WriteString(src)
}

func (se *JavaScriptImpl) Run() (err error) {
	_, err = se.vm.Run(se.sdk + se.script.String())
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

func NewWithContext(ctx core.Context) *JavaScriptImpl {
	se := New()
	se.context = ctx
	se.AddFunc("require", SysRequire(se.context))
	return se
}

func NewWithProcessor(processor core.Processor) *JavaScriptImpl {
	se := New()
	se.context = processor
	se.AddFunc("_http_get", HttpGet)
	se.AddFunc("_http_post", HttpPost)
	se.AddFunc("_http_delete", HttpDelete)
	se.AddFunc("_http_put", HttpPut)
	se.AddFunc("_sys_log", SysLog(se.context.(core.Processor)))
	se.AddFunc("_sys_err", SysError(se.context.(core.Processor)))
	se.AddFunc("require", SysRequire(se.context))
	se.AddFunc("_resp_write", RespWrite(se.context.(core.Processor)))
	se.AddFunc("_db_query", DBQuery(se.context))
	se.AddFunc("_db_queryOne", DBQueryOne(se.context))
	se.AddFunc("_db_save", DBSave(se.context))
	se.AddFunc("_db_exec", DBExec(se.context))
	se.AddFunc("_kv_get", KvGet(se.context))
	se.AddFunc("_kv_set", KvSet(se.context))
	se.AddFunc("_jwt_create", CreateToken)
	se.AddFunc("_jwt_parse", ParseToken)
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
