package services

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zfd81/rock/script"

	"github.com/zfd81/rock/core"

	"github.com/zfd81/rock/errs"

	"github.com/spf13/cast"

	"github.com/zfd81/rooster/util"

	"github.com/zfd81/rock/meta"

	"github.com/zfd81/rock/httpclient"
)

const (
	Regex     = "[A-Za-z0-9_.]+"
	LogFormat = "[LOG] %s "
)

type RockResource struct {
	env       core.Environment
	namespace string            //命名空间 注:不能包含"/"
	se        core.Script       // 脚本引擎
	method    string            // 资源请求方法
	path      string            // 资源原始路径
	regexPath string            // 正则表达式形式路径
	level     int               // 资源级别
	params    []*meta.Parameter // 服务参数
	log       *bytes.Buffer
	resp      *httpclient.Response
	header    http.Header
}

func (r *RockResource) SetEnvironment(env core.Environment) *RockResource {
	r.env = env
	return r
}

func (r *RockResource) GetModule(path string) core.Module {
	return r.env.SelectModule(r.GetNamespace(), path)
}

func (r *RockResource) GetDataSource(name string) core.DB {
	return r.env.SelectDataSource(r.GetNamespace(), name)
}

func (r *RockResource) GetMethod() string {
	return r.method
}

func (r *RockResource) GetPath() string {
	return r.path
}

func (r *RockResource) GetRegexPath() string {
	return r.regexPath
}

func (r *RockResource) GetLevel() int {
	return r.level
}

func (r *RockResource) GetParams() []*meta.Parameter {
	return r.params
}

func (r *RockResource) AddPathParam(param *meta.Parameter) {
	param.Scope = meta.ScopePath
	r.params = append(r.params, param)
}

func (r *RockResource) AddRequestParam(param *meta.Parameter) {
	param.Scope = meta.ScopeRequest
	r.params = append(r.params, param)
}

func (r *RockResource) AddHeaderParam(param *meta.Parameter) {
	param.Scope = meta.ScopeHeader
	r.params = append(r.params, param)
}

func (r *RockResource) GetNamespace() string {
	if r.namespace == "" {
		return meta.DefaultNamespace
	}
	return r.namespace
}

func (r *RockResource) Println(args ...interface{}) error {
	r.log.WriteString(fmt.Sprintf("[INFO] %s ", time.Now().Format("2006-01-02 15:04:05.000")))
	for _, arg := range args {
		r.log.WriteString(cast.ToString(arg))
		r.log.WriteString(" ")
	}
	r.log.WriteString("\n")
	return nil
}

func (r *RockResource) Perror(args ...interface{}) error {
	r.log.WriteString(fmt.Sprintf("[ERROR] %s ", time.Now().Format("2006-01-02 15:04:05.000")))
	for _, arg := range args {
		r.log.WriteString(errs.ErrorStyleFunc(cast.ToString(arg), " "))
	}
	r.log.WriteString("\n")
	return nil
}

func (r *RockResource) SetRespStatus(code int) {
	r.resp.SetStatusCode(code)
}

func (r *RockResource) AddRespHeader(name string, value interface{}) {
	r.resp.AddHeader(name, value)
}

func (r *RockResource) SetRespData(data interface{}) {
	r.resp.SetData(data)
}

func (r *RockResource) Run() (string, *httpclient.Response, error) {
	//添加服务参数
	for _, p := range r.GetParams() {
		r.se.AddVar(p.Name, p.GetValue())
	}

	err := r.se.Run()
	if err != nil {
		r.log.WriteString(fmt.Sprintf(LogFormat, time.Now().Format("2006-01-02 15:04:05.000")))
		r.log.WriteString(errs.ErrorStyleFunc(err))
		r.log.WriteString("\n")
	}
	return r.log.String(), r.resp, err
}

func (r *RockResource) Clear() {
	r.log.Reset()
	r.resp.Clear()
}

func NewResource(serv *meta.Service) *RockResource {
	path := serv.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	res := &RockResource{
		path:   path,
		params: []*meta.Parameter{},
	}
	regexPath, err := util.ReplaceBetween(path, "{", "}", func(i int, s int, e int, c string) (string, error) {
		param, _ := meta.NewParameter(c, "string", meta.ScopePath)
		res.AddPathParam(param)
		return Regex, nil
	})
	if err != nil {
		return nil
	}
	se := script.NewWithProcessor(res)
	se.SetScript(serv.Source)

	res.se = se
	res.method = strings.ToUpper(serv.Method)
	res.regexPath = regexPath
	pathFragments := strings.Split(regexPath, "/")
	res.level = len(pathFragments) - 1
	index := 0
	for i, fragment := range pathFragments {
		if Regex == fragment {
			res.params[index].Index = i
			index++
		}
	}
	for _, p := range serv.Params {
		param := *p
		if param.IsHeaderScope() {
			res.AddHeaderParam(&param)
		} else {
			res.AddRequestParam(&param)
		}
	}
	res.log = new(bytes.Buffer)
	res.resp = &httpclient.Response{
		Header: httpclient.Header{},
	}
	return res
}
