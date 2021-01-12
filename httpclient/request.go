package httpclient

import (
	"net/http"

	"github.com/spf13/cast"
)

type Request struct {
	Request *http.Request
}

func (r *Request) GetHeader(name string) string {
	return r.Request.Header.Get(name)
}

func (r *Request) AddHeader(name string, value interface{}) {
	r.Request.Header.Set(name, cast.ToString(value))
}

func (r *Request) GetPath() string {
	return r.Request.URL.Path
}

func NewRequest(request *http.Request) *Request {
	return &Request{
		Request: request,
	}
}
