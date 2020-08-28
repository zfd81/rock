package meta

import "strings"

type ParamType int

const (
	PathParam = iota
	RequestParam
)

type Parameter struct {
	Name         string      `yaml:"name"` //参数名称
	Comment      string      //参数注解
	DataType     string      `yaml:"dataType"` // 数据类型
	Value        interface{} // 参数值
	DefaultValue interface{} // 参数默认值
	Index        int         `json:"-"` // 参数索引
	Type         ParamType   `json:"-"` //参数类型
}
type Service struct {
	Name   string       `yaml:"name"`
	Path   string       `yaml:"path"`
	Method string       `yaml:"method"`
	Params []*Parameter `yaml:"params"`
	Script string       `yaml:"script"`
}

func FormatServiceName(name string) string {
	if strings.HasPrefix(name, "/") {
		return name[1:]
	}
	return name
}
