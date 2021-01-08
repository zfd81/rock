package httpclient

import (
	"github.com/spf13/cast"
)

type Response struct {
	StatusCode int
	Header     map[string]string
	Content    string
	Data       interface{}
}

func (r *Response) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) AddHeader(name string, value interface{}) {
	r.Header[name] = cast.ToString(value)
}

func (r *Response) SetHeader(header map[string]string) {
	r.Header = header
}

func (r *Response) SetContent(json string) {
	r.Content = json
}

func (r *Response) SetData(data interface{}) {
	r.Data = data
}

func (r *Response) Clear() {
	r.StatusCode = 0
	r.Header = map[string]string{}
	r.Content = ""
	r.Data = nil
}

func NewResponse() *Response {
	return &Response{
		StatusCode: 0,
		Header:     map[string]string{},
		Content:    "",
		Data:       nil,
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}
