package meta

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type KV struct {
	Namespace string //命名空间
	KvsName   string //kvs名称
	Key       string
	DataType  string // 数据类型
	Value     interface{}
	TTL       int64 //租约时间
}

func (k *KV) SetValue(value interface{}) error {
	if _, ok := value.(string); ok {
		k.DataType = DataTypeString
	} else if _, ok := value.(int64); ok {
		k.DataType = DataTypeInteger
	} else if _, ok := value.(bool); ok {
		k.DataType = DataTypeBool
	} else if _, ok := value.(map[string]interface{}); ok {
		k.DataType = DataTypeMap
	} else if _, ok := value.([]string); ok {
		k.DataType = DataTypeStringArray
	} else if _, ok := value.([]int64); ok {
		k.DataType = DataTypeIntegerArray
	} else if _, ok := value.([]map[string]interface{}); ok {
		k.DataType = DataTypeMapArray
	} else {
		return errors.New("Key[" + k.Key + "] data type error")
	}
	k.Value = value
	return nil
}

func (k *KV) EtcdKey() string {
	return KVEtcdKey(k.Namespace, FormatPath(k.KvsName)+"/"+k.Key)
}

func (k *KV) String() (string, error) {
	bytes, err := bson.Marshal(k)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewKv(bytes []byte) (*KV, error) {
	kv := &KV{}
	if err := bson.Unmarshal(bytes, kv); err != nil {
		return nil, err
	}
	return kv, nil
}
