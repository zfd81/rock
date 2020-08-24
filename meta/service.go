package meta

type Parameter struct {
	Name         string `yaml:"name"`
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
