package dai

import (
	"encoding/json"

	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/util/etcd"
)

func CreateDataSource(ds *meta.DataSource) error {
	data, err := json.Marshal(ds)
	if err != nil {
		return errs.NewError(err)
	}
	key := meta.DataSourceKey(ds.Namespace, ds.Name)
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v != nil {
		return errs.New(errs.ErrDsExists)
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func DeleteDataSource(namespace string, name string) (err error) {
	_, err = etcd.Del(meta.DataSourceKey(namespace, name))
	return
}

func ModifyDataSource(ds *meta.DataSource) error {
	data, err := json.Marshal(ds)
	if err != nil {
		return errs.NewError(err)
	}
	key := meta.DataSourceKey(ds.Namespace, ds.Name)
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v == nil {
		return errs.New(errs.ErrDsNotExist)
	}
	_, err = etcd.Put(key, string(data))
	return err
}

func GetDataSource(namespace string, name string) (*meta.DataSource, error) {
	v, err := etcd.Get(meta.DataSourceKey(namespace, name))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	ds := &meta.DataSource{}
	err = json.Unmarshal(v, ds)
	if err != nil {
		return nil, err
	}
	ds.Password = "********"
	return ds, nil
}

func ListDataSource(namespace string, name string) ([]*meta.DataSource, error) {
	dses := make([]*meta.DataSource, 0, 50)
	key := meta.DataSourceKey(namespace, name)
	kvs, err := etcd.GetWithPrefix(key)
	if err == nil {
		for _, kv := range kvs {
			ds := &meta.DataSource{}
			err = json.Unmarshal(kv.Value, ds)
			if err != nil {
				break
			}
			ds.Password = "********"
			dses = append(dses, ds)
		}
	}
	return dses, err
}
