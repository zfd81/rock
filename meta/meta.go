package meta

import (
	"strings"

	"github.com/zfd81/parrot/conf"
)

const (
	MetaDirectory    = "/meta"
	ServiceDirectory = "/serv"
)

var (
	config = conf.GetConfig()
)

func GetMetaRootPath() string {
	return config.Directory + MetaDirectory
}

func GetServiceRootPath() string {
	return GetMetaRootPath() + ServiceDirectory
}

func ServiceKey(method string, path string) string {
	return GetServiceRootPath() + path + config.Meta.NameSeparator + strings.ToLower(method)
}

func ServicePath(key string) (string, string) {
	start := len(GetServiceRootPath())
	end := strings.LastIndex(key, conf.GetConfig().Meta.NameSeparator)
	return key[start:end], key[end+1:]
}

func MetaPath(path string) string {
	start := len(GetMetaRootPath())
	return path[start:]
}

func FormatPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	return path
}
