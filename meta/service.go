package meta

import (
	"encoding/json"
)

type ParamType int

const (
	PathParam = iota
	RequestParam
)

type Service struct {
	Namespace string       //命名空间 注:不能包含"/"
	Name      string       //参数名称
	Path      string       //服务请求路径
	Method    string       //服务请求方法（GET,POST,PUT,DELETE）
	Params    []*Parameter //服务请求参数
	Source    string       //服务执行的脚本
}

func (s *Service) AddParam(name string, dataType string) error {
	param, err := NewParameter(name, dataType)
	if err != nil {
		return err
	}
	s.Params = append(s.Params, param)
	return nil
}

func (s *Service) Key() string {
	return ServiceKey(s.Namespace, s.Method, s.Path)
}

func (s *Service) String() (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewService(jsonstr []byte) (*Service, error) {
	serv := &Service{}
	if err := json.Unmarshal(jsonstr, serv); err != nil {
		return nil, err
	}
	return serv, nil
}
