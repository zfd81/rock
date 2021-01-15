package core

import (
	"github.com/zfd81/rooster/rsql"

	"github.com/zfd81/rooster/types/container"
)

type Function interface {
	Name() string
	Perform(args ...interface{}) (interface{}, error)
}

type Script interface {
	AddVar(name string, value interface{}) error
	GetVar(name string) (interface{}, error)
	GetMlVar(name string) (interface{}, error)
	AddFunc(name string, function interface{}) error
	GetFunc(name string) (Function, error)
	GetMlFunc(name string) (Function, error)
	CallFunc(name string, args ...interface{}) (interface{}, error)
	SetScript(src string)
	AddScript(src string)
	Run() error
}

type Context interface {
	GetNamespace() string
	GetModule(path string) Module
	GetDataSource(name string) DB
}

type Processor interface {
	Context
	Println(args ...interface{}) error
	Perror(args ...interface{}) error
	SetRespStatus(code int)
	AddRespHeader(name string, value interface{})
	SetRespData(data interface{})
}

type DB interface {
	GetNamespace() string
	GetName() string
	QueryMap(query string, arg interface{}) (container.Map, error)
	QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]container.Map, error)
	Query(query string, arg interface{}) (*rsql.Rows, error)
	Exec(query string, arg interface{}) (int64, error)
	Save(arg interface{}, table ...string) (int64, error)
	BatchSave(arg []interface{}, table ...string) (int64, error)
}
