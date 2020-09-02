package dai

import (
	"encoding/json"
	"strings"

	"github.com/zfd81/parrot/core"
	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/parrot/util/etcd"
)

/**
* 服务的key是: /parrot/serv/服务路径.请求方法
**/

func servKey(method string, path string) string {
	return config.Meta.Path +
		config.Meta.ServicePath +
		path + config.Meta.NameSeparator + strings.ToLower(method)
}

func CreateService(serv *meta.Service) error {
	data, err := json.Marshal(serv)
	if err != nil {
		return err
	}
	key := servKey(serv.Method, serv.Path)
	v, err := etcd.Get(key)
	if err != nil {
		return err
	}
	if v != nil {
		return core.ErrServExists
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func DeleteService(method string, path string) (err error) {
	_, err = etcd.Del(servKey(method, path))
	return
}

func ModifyService(serv *meta.Service) error {
	data, err := json.Marshal(serv)
	if err != nil {
		return err
	}
	key := servKey(serv.Method, serv.Path)
	v, err := etcd.Get(key)
	if err != nil {
		return err
	}
	if v == nil {
		return core.ErrServNotExist
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func GetService(method string, path string) (*meta.Service, error) {
	v, err := etcd.Get(servKey(method, path))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	serv := &meta.Service{}
	err = json.Unmarshal(v, serv)
	if err != nil {
		return nil, err
	}
	return serv, nil
}

func ListService(path string) ([]*meta.Service, error) {
	path = config.Meta.Path +
		config.Meta.ServicePath +
		path
	servs := make([]*meta.Service, 0, 50)
	kvs, err := etcd.GetWithPrefix(path)
	if err == nil {
		for _, kv := range kvs {
			serv := &meta.Service{}
			err = json.Unmarshal(kv.Value, serv)
			if err != nil {
				break
			}
			servs = append(servs, serv)
		}
	}
	return servs, err
}