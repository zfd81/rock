package env

import (
	"net/http"
	"strings"

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
