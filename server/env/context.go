package env

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zfd81/parrot/util/etcd"

	"github.com/zfd81/parrot/meta"
)

var (
	getResources    = make(map[int][]Resource) // GET资源映射
	postResources   = make(map[int][]Resource) // POST资源映射
	putResources    = make(map[int][]Resource) // PUT资源映射
	deleteResources = make(map[int][]Resource) // DELETE资源映射
)

func GetResources() map[int][]Resource {
	return getResources
}

func PostResources() map[int][]Resource {
	return postResources
}

func PutResources() map[int][]Resource {
	return putResources
}

func DeleteResources() map[int][]Resource {
	return deleteResources
}

func AddResource(resource Resource) {
	level := resource.GetLevel()
	var resourceMap map[int][]Resource
	if resource.GetMethod() == http.MethodGet {
		resourceMap = getResources
	} else if resource.GetMethod() == http.MethodPost {
		resourceMap = postResources
	} else if resource.GetMethod() == http.MethodPut {
		resourceMap = putResources
	} else if resource.GetMethod() == http.MethodDelete {
		resourceMap = deleteResources
	}
	if resourceMap[level] == nil {
		resourceMap[level] = []Resource{}
	}
	resourceMap[level] = append(resourceMap[level], resource)
}

func RemoveResource(method string, path string) {
	if path != "" || strings.TrimSpace(path) != "" {
		path = meta.FormatPath(path)
		level := len(strings.Split(path, "/")) - 1
		var resourceMap map[int][]Resource
		if strings.ToUpper(method) == http.MethodGet {
			resourceMap = getResources
		} else if strings.ToUpper(method) == http.MethodPost {
			resourceMap = postResources
		} else if strings.ToUpper(method) == http.MethodPut {
			resourceMap = putResources
		} else if strings.ToUpper(method) == http.MethodDelete {
			resourceMap = deleteResources
		}
		resources := resourceMap[level]
		if resources != nil && len(resources) > 0 {
			for i, v := range resources {
				if path == v.GetPath() {
					resourceMap[level] = append(resources[:i], resources[i+1:]...)
					break
				}
			}
		}
	}
}

func InitResources() error {
	kvs, err := etcd.GetWithPrefix(meta.GetServiceRootPath())
	cnt := 0
	if err == nil {
		for _, kv := range kvs {
			serv := &meta.Service{}
			err = json.Unmarshal(kv.Value, serv)
			if err != nil {
				log.Fatal(err)
			}
			res := NewResource(serv)
			AddResource(res)
			path, _ := meta.ServicePath(string(kv.Key))
			fmt.Printf("[INFO] Service %s:%s initialized successfully \n", res.GetMethod(), path)
			cnt++
		}
		fmt.Printf("[INFO] A total of %d services were initialized \n", cnt)
	}
	return err
}

func metaWatcher(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
	full_path := meta.MetaPath(string(key))
	if operType == etcd.CREATE {
		switch {
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			serv := &meta.Service{}
			err := json.Unmarshal(value, serv)
			if err != nil {
				log.Fatal(err)
				return
			}
			res := NewResource(serv)
			AddResource(res)
			path, _ := meta.ServicePath(string(key))
			fmt.Printf("[INFO] Service %s:%s created successfully \n", res.GetMethod(), path)
			break
		}
	} else if operType == etcd.MODIFY {
		switch {
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			serv := &meta.Service{}
			err := json.Unmarshal(value, serv)
			if err != nil {
				log.Fatal(err)
				return
			}
			res := NewResource(serv)
			RemoveResource(serv.Method, serv.Path)
			AddResource(res)
			path, _ := meta.ServicePath(string(key))
			fmt.Printf("[INFO] Service %s:%s modified successfully \n", res.GetMethod(), path)
			break
		}
	} else if operType == etcd.DELETE {
		switch {
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			path, method := meta.ServicePath(string(key))
			RemoveResource(method, path)
			fmt.Printf("[INFO] Service %s:%s deleted successfully \n", strings.ToUpper(method), path)
			break
		}
	}
}

func WatchMeta() {
	etcd.WatchWithPrefix(meta.GetMetaRootPath(), metaWatcher)
}
