package core

import (
	"net/http"

	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
)

type Context interface {
	GetModule(path string) Module
	GetDataSource(name string) DB
	SetHeader(Header http.Header)
	GetHeader() http.Header
}

type Resource interface {
	SetContext(context Context)
	GetContext() Context
	GetNamespace() string
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
