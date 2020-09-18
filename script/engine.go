package script

import (
	"github.com/zfd81/rooster/rsql"

	"github.com/zfd81/rooster/types/container"

	"github.com/robertkrimen/otto"
)

type Function func(call otto.FunctionCall) otto.Value

type ScriptEngine interface {
	AddVar(name string, value interface{}) error
	GetVar(name string) (interface{}, error)
	AddFunc(name string, function Function) error
	SetScript(src string)
	AddScript(src string)
	Run() error
}

type Environment interface {
	GetNamespace() string
	SelectDataSource(name string) DB
}

type Process interface {
	Println(args ...interface{}) error
	Perror(args ...interface{}) error
	SetRespStatus(code int)
	AddRespHeader(name string, value interface{})
	SetRespData(data interface{})
}

type DB interface {
	QueryMap(query string, arg interface{}) (container.Map, error)
	QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]container.Map, error)
	Query(query string, arg interface{}) (*rsql.Rows, error)
	Exec(query string, arg interface{}) (int64, error)
	Save(arg interface{}, table ...string) (int64, error)
	BatchSave(arg []interface{}, table ...string) (int64, error)
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}
