package dai

import (
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/util/etcd"
)

func CreateDataSource(ds *meta.DataSource) error {
	jsonstr, err := ds.String()
	if err != nil {
		return errs.NewError(err)
	}
	key := ds.Key()
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v != nil {
		return errs.New(errs.ErrDsExists)
	}
	_, err = etcd.Put(key, jsonstr)
	return err
}

func DeleteDataSource(ds *meta.DataSource) error {
	_, err := etcd.Del(ds.Key())
	return err
}

func ModifyDataSource(ds *meta.DataSource) error {
	jsonstr, err := ds.String()
	if err != nil {
		return errs.NewError(err)
	}
	key := ds.Key()
	v, err := etcd.Get(key)
	if err != nil {
		return errs.NewError(err)
	}
	if v == nil {
		return errs.New(errs.ErrDsNotExist)
	}
	_, err = etcd.Put(key, jsonstr)
	return err
}

func GetDataSource(namespace string, name string) (*meta.DataSource, error) {
	key := meta.DataSourceKey(namespace, name)
	v, err := etcd.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	ds, err := meta.NewDataSource(v)
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
			ds, err := meta.NewDataSource(kv.Value)
			if err != nil {
				continue
			}
			ds.Password = "********"
			dses = append(dses, ds)
		}
	}
	return dses, err
}
