package services

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/script"
)

type RockModule struct {
	namespace string //命名空间 注:不能包含"/"
	path      string // 模块访问路径
	name      string //模块名称
	source    string //源码
}

func (m *RockModule) GetNamespace() string {
	if m.namespace == "" {
		return meta.DefaultNamespace
	}
	return m.namespace
}

func (m *RockModule) GetPath() string {
	return m.path
}

func (m *RockModule) GetName() string {
	return m.name
}

func (m *RockModule) GetSource() string {
	return m.source
}

func (m *RockModule) GenerateInterceptor() *RockInterceptor {
	se := script.New()
	se.AddScript(script.GetSdk())
	se.AddScript(m.GetSource())
	err := se.Run()
	if err != nil {
		log.Error(err)
		return nil
	}
	paths, err := se.GetMlVar("exports.interceptor.paths")
	if err != nil {
		log.Error(err)
		return nil
	}
	if paths != nil {
		if val, ok := paths.([]string); ok {
			level, _ := se.GetMlVar("exports.interceptor.level")
			interceptor := NewInterceptor(m, val, cast.ToInt(level))
			requestHandler, err := se.GetMlFunc("exports.interceptor.requestHandler")
			if err != nil {
				log.Error(err)
				return nil
			}
			if requestHandler != nil {
				interceptor.SetRequestHandler(requestHandler)
			}
			responseHandler, err := se.GetMlFunc("exports.interceptor.responseHandler")
			if err != nil {
				log.Error(err)
				return nil
			}
			if responseHandler != nil {
				interceptor.SetResponseHandler(responseHandler)
			}
			return interceptor
		}
	}
	return nil
}

func NewModule(serv *meta.Service) *RockModule {
	module := &RockModule{
		namespace: serv.Namespace,
		path:      meta.FormatPath(serv.Path),
		name:      serv.Name,
		source:    serv.Source,
	}
	return module
}
