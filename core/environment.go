package core

import (
	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
)

type Environment interface {
	GetNamespace() string
	AddModule(module Module)
	RemoveModule(namespace string, path string) Module
	SelectModule(namespace string, path string) Module
	GetResourceSet(method string, level int) []Resource
	AddResource(resource Resource)
	RemoveResource(method string, path string)
	SelectResource(method string, path string) Resource
	AddDataSource(ds *meta.DataSource) error
	RemoveDataSource(namespace string, name string) DB
	SelectDataSource(namespace string, name string) DB
}

type Module interface {
	GetNamespace() string
	GetPath() string
	GetName() string
	GetSource() string
}

type Resource interface {
	GetMethod() string
	GetPath() string
	GetRegexPath() string
	GetLevel() int
	GetParams() []*meta.Parameter
	AddPathParam(param *meta.Parameter)
	AddRequestParam(param *meta.Parameter)
	AddHeaderParam(param *meta.Parameter)
	Run() (log string, resp *httpclient.Response, err error)
	Clear()
}
