package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

const (
	DataTypeString       = "STRING"
	DataTypeInteger      = "INT"
	DataTypeBool         = "BOOL"
	DataTypeMap          = "MAP"
	DataTypeStringArray  = "STRING[]"
	DataTypeIntegerArray = "INT[]"
	DataTypeMapArray     = "MAP[]"
)

type Parameter struct {
	Name         string      //参数名称
	Comment      string      `json:"omitempty"` //参数注解
	DataType     string      // 数据类型
	value        interface{} // 参数值
	DefaultValue interface{} `json:"omitempty"` // 参数默认值
	Index        int         `json:"-"`         // 参数索引
	Type         ParamType   `json:"-"`         //参数类型
}

func (p *Parameter) SetValue(value interface{}) error {
	var val interface{}
	var err error
	if strings.ToUpper(p.DataType) == DataTypeString {
		val, err = cast.ToStringE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeInteger {
		val, err = cast.ToIntE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeBool {
		val, err = cast.ToBoolE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeMap {
		val, err = cast.ToStringMapE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeStringArray {
		val, err = cast.ToStringSliceE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeIntegerArray {
		val, err = cast.ToIntSliceE(value)
	} else if strings.ToUpper(p.DataType) == DataTypeMapArray {
		v, ok := value.(string)
		if ok {
			var _val []map[string]interface{}
			err = json.Unmarshal([]byte(v), &_val)
			val = _val
		} else {
			v, ok := value.([]interface{})
			if ok {
				arr := []map[string]interface{}{}
				for _, item := range v {
					m, ok := item.(map[string]interface{})
					if ok {
						arr = append(arr, m)
					} else {
						return fmt.Errorf("unable to cast %#v of type %T to map[]", v, v)
					}
				}
				val = arr
			} else {
				err = fmt.Errorf("unable to cast %#v of type %T to map[]", v, v)
			}
		}
	}
	p.value = val
	return err
}

func (p *Parameter) GetValue() interface{} {
	return p.value
}

func NewParameter(name string, dataType string) (*Parameter, error) {
	if strings.ToUpper(dataType) != DataTypeString &&
		strings.ToUpper(dataType) != DataTypeInteger &&
		strings.ToUpper(dataType) != DataTypeBool &&
		strings.ToUpper(dataType) != DataTypeMap &&
		strings.ToUpper(dataType) != DataTypeStringArray &&
		strings.ToUpper(dataType) != DataTypeIntegerArray &&
		strings.ToUpper(dataType) != DataTypeMapArray {
		return nil, errors.New("Parameter data type error")
	}
	return &Parameter{
		Name:     name,
		DataType: dataType,
	}, nil
}
