package meta

import (
	"strings"

	"github.com/zfd81/rock/conf"
)

const (
	PathSeparator = "/"
	NameSeparator = "."

	DefaultNamespace    = "/__"
	MetaDirectory       = "/meta"
	ServiceDirectory    = "/_serv"
	DataSourceDirectory = "/_ds"
	KVDirectory         = "/_kv"
)

var (
	config = conf.GetConfig()
)

func GetMetaRootPath() string {
	return config.Directory + MetaDirectory
}

//在元数据下的路径
func MetaPath(path string) string {
	start := len(GetMetaRootPath())
	return path[start:]
}

func GetServiceRootPath() string {
	return GetMetaRootPath() + ServiceDirectory
}

func ServiceKey(namespace string, method string, path string) string {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return GetServiceRootPath() + FormatPath(namespace) + PathSeparator + strings.ToLower(method) + FormatPath(path)
}

func ServicePath(key string) (namespace string, method string, path string) {
	strLen := len(key) //字符串长度
	cnt := 0
	position := 0
	for i := 0; i < strLen; i++ {
		char := key[i]
		if char == '/' {
			if ServiceDirectory == key[position:i] {
				cnt++
			} else if cnt == 1 {
				namespace = key[position:i]
				cnt++
			} else if cnt == 2 {
				method = key[position+1 : i]
				cnt++
			} else if cnt == 3 {
				break
			}
			position = i
		}
	}
	path = key[position:]
	return
}

func GetDataSourceRootPath() string {
	return GetMetaRootPath() + DataSourceDirectory
}

func DataSourceKey(namespace string, name string) string {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return GetDataSourceRootPath() + FormatPath(namespace) + FormatPath(name)
}

func DataSourcePath(key string) (namespace string, name string) {
	strLen := len(key) //字符串长度
	cnt := 0
	position := 0
	for i := 0; i < strLen; i++ {
		char := key[i]
		if char == '/' {
			if DataSourceDirectory == key[position:i] {
				cnt++
			} else if cnt == 1 {
				namespace = key[position:i]
				cnt++
			} else if cnt == 2 {
				break
			}
			position = i
		}
	}
	name = key[position:]
	return
}

func GetKVRootPath(namespace string) string {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return GetMetaRootPath() + FormatPath(namespace) + KVDirectory
}

func KVEtcdKey(namespace string, name string) string {
	return GetKVRootPath(namespace) + FormatPath(name)
}

func FormatPath(path string) string {
	if path != "/" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if strings.HasSuffix(path, "/") {
			path = path[0 : len(path)-1]
		}
	}
	return path
}
