package http

import "github.com/spf13/cast"

type Response struct {
	StatusCode int
	Header     map[string]string
	Content    string
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
