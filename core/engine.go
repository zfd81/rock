package core

import (
	"github.com/zfd81/rooster/rsql"

	"github.com/zfd81/rooster/types/container"
)

type Script interface {
	AddVar(name string, value interface{}) error
	GetVar(name string) (interface{}, error)
	AddFunc(name string, function interface{}) error
	CallFunc(name string, args ...interface{}) (interface{}, error)
	SetScript(src string)
	AddScript(src string)
	Run() error
}

type Environment interface {
	GetNamespace() string
	SelectModule(path string) Module
	SelectDataSource(name string) DB
}

type Processor interface {
	Environment
	Println(args ...interface{}) error
	Perror(args ...interface{}) error
	SetRespStatus(code int)
	AddRespHeader(name string, value interface{})
	SetRespData(data interface{})
}

type Module interface {
	GetNamespace() string
	GetPath() string
	GetName() string
	GetSource() string
}

type DB interface {
	QueryMap(query string, arg interface{}) (container.Map, error)
	QueryMapList(query string, arg interface{}, pageNumber int, pageSize int) ([]container.Map, error)
	Query(query string, arg interface{}) (*rsql.Rows, error)
	Exec(query string, arg interface{}) (int64, error)
	Save(arg interface{}, table ...string) (int64, error)
	BatchSave(arg []interface{}, table ...string) (int64, error)
}
