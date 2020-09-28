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

func GetServiceRootPath(namespace string) string {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return GetMetaRootPath() + FormatPath(namespace) + ServiceDirectory
}

func ServiceEtcdKey(namespace string, method string, path string) string {
	return GetServiceRootPath(namespace) + PathSeparator + strings.ToLower(method) + FormatPath(path)
}

func GetDataSourceRootPath(namespace string) string {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return GetMetaRootPath() + FormatPath(namespace) + DataSourceDirectory
}

func DataSourceEtcdKey(namespace string, name string) string {
	return GetDataSourceRootPath(namespace) + FormatPath(name)
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
