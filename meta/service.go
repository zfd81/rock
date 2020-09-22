package meta

import (
	"encoding/json"
)

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
	Name         string      //参数名称
	Comment      string      //参数注解
	DataType     string      // 数据类型
	Value        interface{} // 参数值
	DefaultValue interface{} // 参数默认值
	Index        int         `json:"-"` // 参数索引
	Type         ParamType   `json:"-"` //参数类型
}

type Service struct {
	Namespace string       //命名空间 注:不能包含"/"
	Name      string       //参数名称
	Path      string       //服务请求路径
	Method    string       //服务请求方法（GET,POST,PUT,DELETE）
	Params    []*Parameter //服务请求参数
	Source    string       //服务执行的脚本
}

func (s *Service) AddParam(name string, dataType string) {
	s.Params = append(s.Params, &Parameter{Name: name, DataType: dataType})
}

func NewService(jsonstr []byte) (*Service, error) {
	serv := &Service{}
	if err := json.Unmarshal(jsonstr, serv); err != nil {
		return nil, err
	}
	return serv, nil
}
