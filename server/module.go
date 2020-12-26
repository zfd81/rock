package server

import (
	"github.com/zfd81/rock/meta"
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

func NewModule(serv *meta.Service) *RockModule {
	module := &RockModule{
		namespace: serv.Namespace,
		path:      meta.FormatPath(serv.Path),
		name:      serv.Name,
		source:    serv.Source,
	}
	return module
}
