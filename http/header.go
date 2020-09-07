package http

import "github.com/spf13/cast"

type Header map[string]interface{}

func (h Header) Set(key string, value interface{}) {
	h[key] = value
}

func (h Header) Get(key string) string {
	if v, ok := h[key]; ok {
		return cast.ToString(v)
	} else {
		return ""
	}
}
