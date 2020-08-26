package meta

import "strings"

type Parameter struct {
	Name         string `yaml:"name"`
	Comment      string
	DataType     string `yaml:"dataType"`
	Value        interface{}
	DefaultValue interface{}
	Index        int
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
