package dai

import (
	"encoding/json"

	"github.com/zfd81/parrot/errs"
	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/parrot/util/etcd"
)

/**
* 服务的key是: /parrot/serv/服务路径.请求方法
**/

func CreateService(serv *meta.Service) error {
	data, err := json.Marshal(serv)
	if err != nil {
		return errs.NewError(err)
	}
	key := meta.ServiceKey(serv.Namespace, serv.Method, serv.Path)
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v != nil {
		return errs.New(errs.ErrServExists)
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func DeleteService(serv *meta.Service) (err error) {
	_, err = etcd.Del(meta.ServiceKey(serv.Namespace, serv.Method, serv.Path))
	return
}

func ModifyService(serv *meta.Service) error {
	data, err := json.Marshal(serv)
	if err != nil {
		return errs.NewError(err)
	}
	key := meta.ServiceKey(serv.Namespace, serv.Method, serv.Path)
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v == nil {
		return errs.New(errs.ErrServNotExist)
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func GetService(namespace string, method string, path string) (*meta.Service, error) {
	v, err := etcd.Get(meta.ServiceKey(namespace, method, path))
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

func ListService(namespace string, path string) ([]*meta.Service, error) {
	servs := make([]*meta.Service, 0, 50)

	//查询GET服务
	key := meta.ServiceKey(namespace, "get", path)
	kvs, err := etcd.GetWithPrefix(key)
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

	//查询POST服务
	key = meta.ServiceKey(namespace, "post", path)
	kvs, err = etcd.GetWithPrefix(key)
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

	//查询PUT服务
	key = meta.ServiceKey(namespace, "put", path)
	kvs, err = etcd.GetWithPrefix(key)
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

	//查询DELETE服务
	key = meta.ServiceKey(namespace, "delete", path)
	kvs, err = etcd.GetWithPrefix(key)
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
