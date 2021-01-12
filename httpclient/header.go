package httpclient

import "github.com/spf13/cast"

type Header map[string]string

func (h Header) Set(key string, value interface{}) {
	h[key] = cast.ToString(value)
}

func (h Header) Get(key string) string {
	if v, ok := h[key]; ok {
		return v
	} else {
		return ""
	}
}

func (h Header) Has(key string) bool {
	_, ok := h[key]
	return ok
}

func (h Header) Del(key string) {
	if _, ok := h[key]; ok {
		delete(h, key)
	}
}

func (h Header) Clone() Header {
	if h == nil {
		return nil
	}
	cloneh := Header{}
	for k, v := range h {
		cloneh.Set(k, v)
	}
	return cloneh
}

func (h Header) Clear() Header {
	h = map[string]string{}
	return h
}
