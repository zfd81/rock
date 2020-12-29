package otto

import (
	"strings"

	"github.com/zfd81/rock/core"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/meta"

	"github.com/zfd81/rock/meta/dai"
)

func KvGet(process core.Processor) func(name, key string) interface{} {
	return func(name, key string) interface{} {
		name = strings.TrimSpace(name) //获取kvs名称
		if name == "" {
			throwException("KVS name cannot be empty")
		}
		key = strings.TrimSpace(key) //获取kvs名称
		if key == "" {
			throwException("KV key cannot be empty")
		}
		kv, err := dai.GetKV(process.GetNamespace(), meta.FormatPath(name)+"/"+key)
		if err != nil {
			throwException(err.Error())
		}
		if kv == nil {
			return nil
		}
		return kv.Value
	}
}

func KvSet(process core.Processor) func(name, key string, value interface{}, ttl int64) {
	return func(name, key string, value interface{}, ttl int64) {
		name = strings.TrimSpace(name) //获取kvs名称
		if name == "" {
			throwException("KVS name cannot be empty")
		}
		key = strings.TrimSpace(key) //获取kvs名称
		if key == "" {
			throwException("KV key cannot be empty")
		}
		if value == nil {
			throwException("KV value cannot be empty")
		}
		if ttl < 1 {
			ttl = conf.GetConfig().KVTTL
		}
		kv := &meta.KV{
			Namespace: process.GetNamespace(),
			KvsName:   name,
			Key:       key,
			TTL:       ttl,
		}
		if err := kv.SetValue(value); err != nil {
			throwException(err.Error())
		}
		if err := dai.SetKV(kv); err != nil {
			throwException(err.Error())
		}
	}
}
