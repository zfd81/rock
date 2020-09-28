package dai

import (
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/util/etcd"
)

func SetKV(kv *meta.KV) error {
	str, err := kv.String()
	if err != nil {
		return errs.NewError(err)
	}
	key := kv.EtcdKey()
	_, err = etcd.PutWithTTL(key, str, kv.TTL)
	return err
}

func GetKV(namespace string, name string) (*meta.KV, error) {
	key := meta.KVEtcdKey(namespace, name)
	v, err := etcd.Get(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	kv, err := meta.NewKv(v)
	if err != nil {
		return nil, err
	}
	return kv, nil
}
