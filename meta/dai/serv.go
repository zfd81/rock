package dai

import (
	"encoding/json"
	"errors"

	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/parrot/util/etcd"
)

func CreateService(serv *meta.Service) (int, error) {
	data, err := json.Marshal(serv)
	if err != nil {
		return 0, err
	}
	key := servKey(serv.Name)
	v, err := etcd.Get(key)
	if v != nil {
		return -1, errors.New("The service already exists")
	}
	_, err = etcd.Put(servKey(serv.Name), string(data))
	return 1, err
}

func DeleteService(name string) error {
	_, err := etcd.Del(servKey(name))
	if err != nil {
		return err
	}
	return nil
}

func ModifyService(serv *meta.Service) error {
	data, err := json.Marshal(serv)
	if err != nil {
		return err
	}
	_, err = etcd.Put(servKey(serv.Name), string(data))
	return err
}

func GetService(name string) (*meta.Service, error) {
	v, err := etcd.Get(servKey(name))
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

func ListService(name string) ([]*meta.Service, error) {
	path := config.Meta.RootDirectory +
		config.Meta.ServiceDirectory +
		config.Meta.PathSeparator +
		name
	servs := make([]*meta.Service, 0, 50)
	kvs, err := etcd.GetWithPrefix(path)
	if err == nil {
		for _, kv := range kvs {
			serv := &meta.Service{}
			err = json.Unmarshal(kv.Value, serv)
			if err != nil {
				return servs, err
			}
			servs = append(servs, serv)
		}
	}
	return servs, err
}

func servKey(name string) string {
	return config.Meta.RootDirectory +
		config.Meta.ServiceDirectory +
		config.Meta.PathSeparator +
		name +
		config.Meta.ServiceSuffix
}
