package env

import (
	"bytes"
	"fmt"
	"time"

	"github.com/zfd81/parrot/script"

	"github.com/zfd81/parrot/meta"

	"github.com/spf13/cast"

	"github.com/zfd81/parrot/http"
)

const (
	LogFormat = "[LOG] %s "
)

type Instance struct {
	se   script.ScriptEngine
	log  *bytes.Buffer
	resp http.Response
}

func (i *Instance) SetParam(name string, value interface{}) error {
	return i.se.AddVar(name, value)
}

func (i *Instance) Println(args ...interface{}) error {
	i.log.WriteString(fmt.Sprintf(LogFormat, time.Now().Format("2006-01-02 15:04:05.000")))
	for _, arg := range args {
		i.log.WriteString(cast.ToString(arg))
	}
	i.log.WriteString("\n")
	return nil
}

func (i *Instance) SetRespStatus(code int) {
	i.resp.SetStatusCode(code)
}

func (i *Instance) GetRespStatus() int {
	return i.resp.StatusCode
}

func (i *Instance) AddRespHeader(name string, value interface{}) {
	i.resp.AddHeader(name, value)
}

func (i *Instance) SetRespContent(json string) {
	i.resp.SetContent(json)
}

func (i *Instance) SetRespData(data interface{}) {
	i.resp.SetData(data)
}

func (i *Instance) GetRespContent() string {
	return i.resp.Content
}

func (i *Instance) Run() (string, error) {
	err := i.se.Run()
	log := i.log.String()
	i.log.Reset()
	return log, err
}

func New(serv *meta.Service) *Instance {
	se := script.New()
	ins := &Instance{
		se:   se,
		log:  new(bytes.Buffer),
		resp: http.Response{},
	}
	se.SetScript(serv.Script)
	se.AddFunc("_sys_log", script.SysLog(ins))
	se.AddFunc("_resp_write", script.RespWrite(ins))
	se.AddFunc("_http_get", script.HttpGet)
	se.AddFunc("_http_post", script.HttpPost)
	se.AddFunc("_http_delete", script.HttpDelete)
	se.AddFunc("_http_put", script.HttpPut)
	return ins
}
