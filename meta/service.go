package meta

type ParamType int

const (
	PathParam = iota
	RequestParam

	DataTypeString  = "STRING"
	DataTypeInteger = "INT"
	DataTypeBool    = "BOOL"
	DataTypeMap     = "MAP"
	DataTypeArray   = "ARR"
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
	Namespace string       `yaml:"namespace"` //命名空间 注:不能包含"/"
	Name      string       `yaml:"name"`      //参数名称
	Path      string       `yaml:"path"`      //服务请求路径
	Method    string       `yaml:"method"`    //服务请求方法（GET,POST,PUT,DELETE）
	Params    []*Parameter `yaml:"params"`
	Script    string       `yaml:"script"`
}

func (s *Service) AddParam(name string, dataType string) {
	s.Params = append(s.Params, &Parameter{Name: name, DataType: dataType})
}
