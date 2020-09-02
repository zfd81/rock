package env

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zfd81/parrot/conf"

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

func LoadMeta() error {
	kvs, err := etcd.GetWithPrefix(conf.GetConfig().Meta.Path + conf.GetConfig().Meta.ServicePath)
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
			fmt.Printf("[INFO] Service %s:%s initialized successfully \n", res.GetMethod(), servPath(string(kv.Key)))
			cnt++
		}
		fmt.Printf("[INFO] A total of %d services were initialized \n", cnt)
	}
	return err
}

func servPath(path string) string {
	start := len(conf.GetConfig().Meta.Path + conf.GetConfig().Meta.ServicePath)
	end := strings.LastIndex(path, conf.GetConfig().Meta.NameSeparator)
	return path[start:end]
}
