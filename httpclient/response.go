package httpclient

type Response struct {
	StatusCode int
	Header     Header
	Content    string
	Data       interface{}
}

func (r *Response) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) GetHeader(name string) string {
	return r.Header.Get(name)
}

func (r *Response) AddHeader(name string, value interface{}) {
	r.Header.Set(name, value)
}

func (r *Response) SetContent(json string) {
	r.Content = json
}

func (r *Response) SetData(data interface{}) {
	r.Data = data
}

func (r *Response) Clear() {
	r.StatusCode = 0
	r.Header.Clear()
	r.Content = ""
	r.Data = nil
}

func NewResponse() *Response {
	return &Response{
		StatusCode: 0,
		Header:     Header{},
		Content:    "",
		Data:       nil,
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}
