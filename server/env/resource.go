package env

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/zfd81/rooster/util"

	"github.com/zfd81/parrot/meta"

	"github.com/zfd81/parrot/http"
	"github.com/zfd81/parrot/script"
)

const (
	Regex = "[A-Za-z0-9_.]+"
)

type Resource interface {
	GetMethod() string
	GetPath() string
	GetRegexPath() string
	GetResourcePath() string
	SetResourcePath(resourcePath string)
	GetLevel() int
	Run() (string, error)
}

type ParrotResource struct {
	se            script.ScriptEngine // 脚本引擎
	method        string              // 资源请求方法
	path          string              // 资源原始路径
	regexPath     string              // 正则表达式形式路径
	resourcePath  string              // 资源路径
	level         int                 // 资源级别
	PathParams    []*meta.Parameter   // 路径参数
	RequestParams []*meta.Parameter   //请求参数
	log           *bytes.Buffer
	resp          http.Response
}

func (r *ParrotResource) GetMethod() string {
	return r.method
}

func (r *ParrotResource) GetPath() string {
	return r.path
}

func (r *ParrotResource) GetRegexPath() string {
	return r.regexPath
}

func (r *ParrotResource) GetResourcePath() string {
	return r.resourcePath
}

func (r *ParrotResource) SetResourcePath(resourcePath string) {
	r.resourcePath = resourcePath
}

func (r *ParrotResource) GetLevel() int {
	return r.level
}

func (r *ParrotResource) Println(args ...interface{}) error {
	r.log.WriteString(fmt.Sprintf(LogFormat, time.Now().Format("2006-01-02 15:04:05.000")))
	for _, arg := range args {
		r.log.WriteString(cast.ToString(arg))
	}
	r.log.WriteString("\n")
	return nil
}

func (r *ParrotResource) SetRespStatus(code int) {
	r.resp.SetStatusCode(code)
}

func (r *ParrotResource) AddRespHeader(name string, value interface{}) {
	r.resp.AddHeader(name, value)
}

func (r *ParrotResource) SetRespContent(json string) {
	r.resp.SetContent(json)
}

func (r *ParrotResource) Run() (string, error) {
	err := r.se.Run()
	log := r.log.String()
	r.log.Reset()
	return log, err
}

func NewResource(serv *meta.Service) Resource {
	path := serv.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	res := &ParrotResource{
		path:       path,
		PathParams: []*meta.Parameter{},
	}
	regexPath, err := util.ReplaceBetween(path, "{", "}", func(i int, s int, e int, c string) (string, error) {
		param := &meta.Parameter{
			Name:     c,
			DataType: "string",
		}
		res.PathParams = append(res.PathParams, param)
		return Regex, nil
	})
	if err != nil {
		return nil
	}
	se := script.New()
	se.SetScript(serv.Script)
	se.AddFunc("_sys_log", script.SysLog(res))
	se.AddFunc("_resp_write", script.RespWrite(res))
	se.AddFunc("_http_get", script.HttpGet)
	se.AddFunc("_http_post", script.HttpPost)
	se.AddFunc("_http_delete", script.HttpDelete)
	se.AddFunc("_http_put", script.HttpPut)
	res.se = se
	res.method = strings.ToUpper(serv.Method)
	res.regexPath = regexPath
	pathFragments := strings.Split(regexPath, "/")
	index := 0
	for i, fragment := range pathFragments {
		if Regex == fragment {
			res.PathParams[index].Index = i
			index++
		}
	}
	res.level = len(pathFragments) - 1
	res.log = new(bytes.Buffer)
	res.resp = http.Response{}
	return res
}
