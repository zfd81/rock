package httpclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

var client = &http.Client{}

func Do(req *http.Request, header map[string]string) (*http.Response, error) {
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return client.Do(req)
}

func Get(url string, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return Do(req, header)
}

func Post(url, contentType string, body io.Reader, header map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return Do(req, header)
}

func PostForm(url string, data url.Values, header map[string]string) (resp *http.Response, err error) {
	return Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), header)
}
