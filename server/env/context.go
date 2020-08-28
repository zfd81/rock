package env

import "net/http"

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
	resourceMap := GetResources()
	if resource.GetMethod() == http.MethodPost {
		resourceMap = PostResources()
	} else if resource.GetMethod() == http.MethodPut {
		resourceMap = PutResources()
	} else if resource.GetMethod() == http.MethodDelete {
		resourceMap = DeleteResources()
	}
	if resourceMap[level] == nil {
		resourceMap[level] = []Resource{}
	}
	resourceMap[level] = append(resourceMap[level], resource)
}
