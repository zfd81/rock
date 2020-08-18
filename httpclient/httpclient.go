package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type Response struct {
	Status     string
	StatusCode int
	Header     map[string]string
	Content    string
}

func (r *Response) wrap(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r.Status = resp.Status
	r.StatusCode = resp.StatusCode
	r.Header = map[string]string{}
	for k, v := range resp.Header {
		r.Header[k] = v[0]
	}
	r.Content = string(body)
	return nil
}

type HttpClient struct {
	client  *http.Client
	Timeout time.Duration
}

func (hc *HttpClient) do(req *http.Request, header map[string]string) (*http.Response, error) {
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return hc.client.Do(req)
}

func (hc *HttpClient) Get(url string, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return hc.do(req, header)
}

func (hc *HttpClient) Post(url, contentType string, body io.Reader, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return hc.do(req, header)
}

func (hc *HttpClient) PostForm(url string, data url.Values, header map[string]string) (resp *http.Response, err error) {
	return hc.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), header)
}

func (hc *HttpClient) Put(url, contentType string, body io.Reader, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return hc.do(req, header)
}

func (hc *HttpClient) Delete(url string, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return hc.do(req, header)
}

var client = HttpClient{
	client: &http.Client{},
}

func Get(url string, header map[string]string) *Response {
	response := &Response{StatusCode: 500}
	resp, err := client.Get(url, header)
	if err != nil {
		log.Println(err)
		return response
	}

	defer func() {
	_:
		resp.Body.Close()
	}()

	err = response.wrap(resp)
	if err != nil {
		log.Println(err)
	}
	return response
}

func Post(url string, data map[string]interface{}, header map[string]string) *Response {
	response := &Response{StatusCode: 500}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Post(url, "application/json;charset=UTF-8", bytes.NewBuffer(jsonStr), header)
	if err != nil {
		log.Println(err)
		return response
	}

	defer func() {
	_:
		resp.Body.Close()
	}()

	err = response.wrap(resp)
	if err != nil {
		log.Println(err)
	}
	return response
}

func PostForm(respUrl string, data map[string]interface{}, header map[string]string) *Response {
	response := &Response{StatusCode: 500}
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values.Add(k, cast.ToString(v))
		}
	}
	resp, err := client.PostForm(respUrl, values, header)
	if err != nil {
		log.Println(err)
		return response
	}

	defer func() {
	_:
		resp.Body.Close()
	}()

	err = response.wrap(resp)
	if err != nil {
		log.Println(err)
	}
	return response
}

func Put(url string, data map[string]interface{}, header map[string]string) *Response {
	response := &Response{StatusCode: 500}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return response
	}

	resp, err := client.Put(url, "application/json;charset=UTF-8", bytes.NewBuffer(jsonStr), header)
	if err != nil {
		log.Println(err)
		return response
	}

	defer func() {
	_:
		resp.Body.Close()
	}()

	err = response.wrap(resp)
	if err != nil {
		log.Println(err)
	}
	return response
}

func Delete(url string, header map[string]string) *Response {
	response := &Response{StatusCode: 500}
	resp, err := client.Delete(url, header)
	if err != nil {
		log.Println(err)
		return response
	}

	defer func() {
	_:
		resp.Body.Close()
	}()

	err = response.wrap(resp)
	if err != nil {
		log.Println(err)
	}
	return response
}
