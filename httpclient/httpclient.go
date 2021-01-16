package httpclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zfd81/rock/conf"

	"github.com/spf13/cast"
	"github.com/zfd81/rooster/util"
)

type HttpClient struct {
	client  *http.Client
	Timeout time.Duration
}

func (hc *HttpClient) do(req *http.Request, header Header) (*Response, error) {
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, cast.ToString(v))
		}
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Header:     Header{},
		Content:    string(body),
	}

	for k, v := range resp.Header {
		response.Header[k] = v[0]
	}

	return response, nil
}

func (hc *HttpClient) Get(url string, data map[string]interface{}, header Header) (resp *Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if data != nil && len(data) > 0 {
		q := req.URL.Query()
		for k, v := range data {
			q.Add(k, cast.ToString(v))
		}
		req.URL.RawQuery = q.Encode()
	}
	return hc.do(req, header)
}

func (hc *HttpClient) Post(url, contentType string, data interface{}, header Header) (resp *Response, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return hc.do(req, header)
}

func (hc *HttpClient) PostForm(url string, data url.Values, header Header) (resp *Response, err error) {
	return hc.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), header)
}

func (hc *HttpClient) Put(url string, data interface{}, header Header) (resp *Response, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return hc.do(req, header)
}

func (hc *HttpClient) Delete(url string, data interface{}, header Header) (resp *Response, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return hc.do(req, header)
}

func New(second time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: time.Second * second,
		},
		Timeout: second,
	}
}

var client = New(conf.GetConfig().HttpClientTimeout)

func Get(url string, data map[string]interface{}, header Header) *Response {
	response := &Response{StatusCode: 500}

	url, err := wrapPath(url, data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Get(url, data, header)
	if err != nil {
		log.Println(err)
		return response
	}

	return resp
}

func Post(url string, data map[string]interface{}, header Header) *Response {
	response := &Response{StatusCode: 500}

	url, err := wrapPath(url, data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Post(url, "application/json;charset=UTF-8", data, header)
	if err != nil {
		log.Println(err)
		return response
	}

	return resp
}

func PostForm(reqUrl string, data map[string]interface{}, header Header) *Response {
	response := &Response{StatusCode: 500}

	reqUrl, err := wrapPath(reqUrl, data)
	if err != nil {
		log.Println(err)
		return response
	}

	values := url.Values{}
	if data != nil && len(data) != 0 {
		for k, v := range data {
			values.Add(k, cast.ToString(v))
		}
	}
	resp, err := client.PostForm(reqUrl, values, header)
	if err != nil {
		log.Println(err)
		return response
	}

	return resp
}

func Put(url string, data map[string]interface{}, header Header) *Response {
	response := &Response{StatusCode: 500}

	url, err := wrapPath(url, data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Put(url, data, header)
	if err != nil {
		log.Println(err)
		return response
	}

	return resp
}

func Delete(url string, data map[string]interface{}, header Header) *Response {
	response := &Response{StatusCode: 500}

	url, err := wrapPath(url, data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Delete(url, data, header)
	if err != nil {
		log.Println(err)
		return response
	}

	return resp
}

func wrapPath(url string, param map[string]interface{}) (string, error) {
	if param != nil {
		return util.ReplaceBetween(url, "{", "}", func(i int, s int, e int, c string) (string, error) {
			key := strings.TrimSpace(c)
			value, found := param[key]
			if found {
				delete(param, key)
				return cast.ToString(value), nil
			}
			return "", nil
		})
	}
	return url, nil
}
