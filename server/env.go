package server

import (
	"regexp"
	"strings"

	"github.com/zfd81/rock/server/services"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/rock/core"

	"github.com/pkg/errors"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
)

type RockEnvironment struct {
	modules   map[string]core.Module
	resources map[string]map[int][]core.Resource
	dbs       map[string]core.DB
}

func (e *RockEnvironment) GetNamespace() string {
	return config.Namespace
}

func (e *RockEnvironment) AddModule(module core.Module) {
	key := meta.FormatPath(module.GetNamespace()) + meta.FormatPath(module.GetPath())
	e.modules[key] = module
}

func (e *RockEnvironment) RemoveModule(namespace string, path string) core.Module {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(path)
	value, found := e.modules[key]
	if found {
		delete(e.modules, key)
		return value
	}
	return nil
}

func (e *RockEnvironment) SelectModule(namespace string, path string) core.Module {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(path)
	value, found := e.modules[key]
	if found {
		return value
	}
	return nil
}

func (e *RockEnvironment) GetResourceSet(method string, level int) []core.Resource {
	value, found := e.resources[method][level]
	if found {
		return value
	}
	return nil
}

func (e *RockEnvironment) AddResource(resource core.Resource) {
	method := resource.GetMethod()
	level := resource.GetLevel()
	rs := e.GetResourceSet(method, level)
	if rs == nil {
		rs = []core.Resource{}
	}
	e.resources[method][level] = append(rs, resource)
}

func (e *RockEnvironment) RemoveResource(method string, path string) {
	if path != "" || strings.TrimSpace(path) != "" {
		path = meta.FormatPath(path)
		method = strings.ToUpper(method)
		level := len(strings.Split(path, "/")) - 1
		rs := e.GetResourceSet(method, level)
		if rs != nil && len(rs) > 0 {
			for i, v := range rs {
				if path == v.GetPath() {
					e.resources[method][level] = append(rs[:i], rs[i+1:]...)
					break
				}
			}
		}
	}
}

func (e *RockEnvironment) SelectResource(method string, path string) core.Resource {
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	level := len(strings.Split(path, "/")) - 1
	rs := e.GetResourceSet(method, level)
	if rs != nil {
		for _, resource := range rs {
			pattern, err := regexp.Compile(resource.GetRegexPath())
			if err != nil {
				log.Println(errors.WithStack(err))
				return nil
			}
			if pattern.MatchString(path) {
				pathFragments := strings.Split(path, "/")
				for _, param := range resource.GetParams() {
					if param.IsPathScope() {
						param.SetValue(pathFragments[param.Index])
					}
				}
				resource.Clear()
				return resource
			}
		}
	}
	return nil
}

func (e *RockEnvironment) AddDataSource(ds *meta.DataSource) error {
	db, err := core.NewDB(ds)
	if err == nil {
		e.dbs[meta.FormatPath(db.GetNamespace())+meta.FormatPath(db.Name)] = db
	}
	return err
}

func (e *RockEnvironment) RemoveDataSource(namespace string, name string) core.DB {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(name)
	value, found := e.dbs[key]
	if found {
		delete(e.dbs, key)
		return value
	}
	return nil
}

func (e *RockEnvironment) SelectDataSource(namespace string, name string) core.DB {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(name)
	//value, found := dbs[key]
	//if found {
	//	return value
	//}
	return e.dbs[key]
}

var (
	config = conf.GetConfig()
	env    = &RockEnvironment{
		modules: map[string]core.Module{}, //模块映射
		resources: map[string]map[int][]core.Resource{
			httpclient.MethodGet:    map[int][]core.Resource{}, // GET资源映射
			httpclient.MethodPost:   map[int][]core.Resource{}, // POST资源映射
			httpclient.MethodPut:    map[int][]core.Resource{}, // PUT资源映射
			httpclient.MethodDelete: map[int][]core.Resource{}, // DELETE资源映射
		},
		dbs: map[string]core.DB{},
	}
	interceptorChain = services.InterceptorChain{} //拦截器链
)

func GetEnvironment() core.Environment {
	return env
}

func GetInterceptorChain() services.InterceptorChain {
	return interceptorChain
}

func AddInterceptor(interceptor *services.RockInterceptor) {
	interceptorChain.Add(interceptor)
}

func RemoveInterceptor(path string) *services.RockInterceptor {
	return interceptorChain.Remove(meta.FormatPath(path))
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
