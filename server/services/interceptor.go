package services

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cast"

	"github.com/zfd81/rock/core"
)

type Handler interface {
	SetNextHandler(handler *Handler)
	GetNextHandler() *Handler
}

type RockInterceptor struct {
	*RockModule
	group           string        //组
	paths           []string      //拦截路径
	level           int           //拦截级别
	requestHandler  core.Function //请求拦截器
	responseHandler core.Function //响应拦截器
}

func (i *RockInterceptor) GetGroup() string {
	return i.group
}

func (i *RockInterceptor) GetPaths() []string {
	return i.paths
}

func (i *RockInterceptor) GetLevel() []string {
	return i.paths
}

func (i *RockInterceptor) AddPath(path string) {
	if path != "" {
		if strings.Index(path, "**") != -1 {
			splitted := strings.SplitN(path, "**", 2)
			path = "^" + splitted[0]
			if splitted[1] == "" {
				path = path + "[A-Za-z0-9_./]#"
			} else {
				path = path + "[A-Za-z0-9_./]+" + splitted[1]
			}
			path = path + "$"
		}
		if strings.Index(path, "*") != -1 {
			splitted := strings.SplitN(path, "*", 2)
			path = splitted[0]
			if !strings.HasPrefix(path, "^") {
				path = "^" + path
			}
			if splitted[1] == "" {
				path = path + "[A-Za-z0-9_.]*"
			} else {
				path = path + "[A-Za-z0-9_.]+" + splitted[1]
			}
			if !strings.HasSuffix(path, "$") {
				path = path + "$"
			}

		}
		path = strings.ReplaceAll(path, "#", "*")
		i.paths = append(i.paths, path)
	}
}

func (i *RockInterceptor) Matches(path string) bool {
	for _, v := range i.paths {
		r, _ := regexp.Compile(v)
		if r.MatchString(path) {
			return true
		}
	}
	return false
}

func (i *RockInterceptor) SetRequestHandler(function core.Function) {
	i.requestHandler = function
}

func (i *RockInterceptor) Request(request *http.Request, response *http.Response) (bool, error) {
	if i.requestHandler != nil {
		val, err := i.requestHandler.Perform(request, response)
		return cast.ToBool(val), err
	}
	return true, nil
}

func (i *RockInterceptor) SetResponseHandler(function core.Function) {
	i.responseHandler = function
}

func (i *RockInterceptor) Response(request *http.Request, response *http.Response) (bool, error) {
	if i.responseHandler != nil {
		val, err := i.responseHandler.Perform(request, response)
		return cast.ToBool(val), err
	}
	return true, nil
}

func NewInterceptor(module *RockModule, paths []string, level int) *RockInterceptor {
	i := &RockInterceptor{
		RockModule: module,
		paths:      paths,
		level:      level,
	}
	return i
}

type InterceptorChain []*RockInterceptor

func (c InterceptorChain) Len() int {
	return len(c)
}

func (c InterceptorChain) Less(i, j int) bool {
	return c[i].level < c[j].level
}

func (c InterceptorChain) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *InterceptorChain) Add(interceptor *RockInterceptor) *InterceptorChain {
	*c = append(*c, interceptor)
	sort.Sort(c)
	return c
}

func (c *InterceptorChain) Remove(path string) *InterceptorChain {
	for i, v := range *c {
		if v.path == path {
			*c = append((*c)[:i], (*c)[i+1:]...)
			break
		}
	}
	return c
}

func (c *InterceptorChain) Modify(interceptor *RockInterceptor) *InterceptorChain {
	for i, v := range *c {
		if v.path == interceptor.path {
			(*c)[i] = interceptor
			break
		}
	}
	sort.Sort(c)
	return c
}
