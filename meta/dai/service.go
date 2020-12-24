package dai

import (
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/util/etcd"
)

func CreateService(serv *meta.Service) error {
	jsonstr, err := serv.String()
	if err != nil {
		return errs.NewError(err)
	}
	key := serv.EtcdKey()
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v != nil {
		return errs.New(errs.ErrServExists)
	}
	_, err = etcd.Put(key, jsonstr)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func DeleteService(serv *meta.Service) error {
	key := serv.EtcdKey()
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v == nil {
		return errs.New(errs.ErrServNotExist)
	}
	_, err = etcd.Del(key)
	return err
}

func ModifyService(serv *meta.Service) error {
	jsonstr, err := serv.String()
	if err != nil {
		return errs.NewError(err)
	}
	key := serv.EtcdKey()
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v == nil {
		return errs.New(errs.ErrServNotExist)
	}
	_, err = etcd.Put(key, jsonstr)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func GetService(namespace string, method string, path string) (*meta.Service, error) {
	key := meta.ServiceEtcdKey(namespace, method, path)
	v, err := etcd.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	serv, err := meta.NewService(v)
	if err != nil {
		return nil, err
	}
	return serv, nil
}

func ListService(namespace string, path string) ([]*meta.Service, error) {
	servs := make([]*meta.Service, 0, 50)

	//查询GET服务
	key := meta.ServiceEtcdKey(namespace, "get", path)
	kvs, err := etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				continue
			}
			servs = append(servs, serv)
		}
	}

	//查询POST服务
	key = meta.ServiceEtcdKey(namespace, "post", path)
	kvs, err = etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				continue
			}
			servs = append(servs, serv)
		}
	}

	//查询PUT服务
	key = meta.ServiceEtcdKey(namespace, "put", path)
	kvs, err = etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				continue
			}
			servs = append(servs, serv)
		}
	}

	//查询DELETE服务
	key = meta.ServiceEtcdKey(namespace, "delete", path)
	kvs, err = etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				continue
			}
			servs = append(servs, serv)
		}
	}

	//查询Module服务
	key = meta.ServiceEtcdKey(namespace, "local", path)
	kvs, err = etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				continue
			}
			servs = append(servs, serv)
		}
	}
	return servs, err
}
