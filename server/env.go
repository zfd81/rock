package server

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/zfd81/rock/script"

	"github.com/zfd81/rock/server/services"

	"github.com/spf13/cast"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/rock/errs"

	"github.com/zfd81/rock/core"

	"github.com/pkg/errors"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/util/etcd"

	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
)

type ResourceContext struct {
	namespace string
	header    http.Header
}

func (c *ResourceContext) GetModule(path string) core.Module {
	return SelectModule(c.namespace, path)
}

func (c *ResourceContext) GetDataSource(name string) core.DB {
	return SelectDataSource(c.namespace, name)
}

func (c *ResourceContext) SetHeader(header http.Header) {
	c.header = header
}

func (c *ResourceContext) GetHeader() http.Header {
	return c.header
}

func NewContext(namespace string) *ResourceContext {
	return &ResourceContext{
		namespace: namespace,
	}
}

var (
	config = conf.GetConfig()

	interceptorChain = services.InterceptorChain{} //拦截器链
	modules          = map[string]core.Module{}    //模块映射
	resources        = map[string]map[int][]core.Resource{
		httpclient.MethodGet:    map[int][]core.Resource{}, // GET资源映射
		httpclient.MethodPost:   map[int][]core.Resource{}, // POST资源映射
		httpclient.MethodPut:    map[int][]core.Resource{}, // PUT资源映射
		httpclient.MethodDelete: map[int][]core.Resource{}, // DELETE资源映射
	}
	dbs = map[string]*core.RockDB{}
)

func AddInterceptor(interceptor *services.RockInterceptor) {
	interceptorChain.Add(interceptor)
}

func RemoveInterceptor(path string) *services.RockInterceptor {
	return interceptorChain.Remove(meta.FormatPath(path))
}

func ModifyInterceptor(module core.Module) {
}

func GetInterceptorChain() services.InterceptorChain {
	return interceptorChain
}

func AddModule(module core.Module) {
	key := meta.FormatPath(module.GetNamespace()) + meta.FormatPath(module.GetPath())
	modules[key] = module
}

func RemoveModule(namespace string, path string) core.Module {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(path)
	value, found := modules[key]
	if found {
		delete(modules, key)
		return value
	}
	RemoveInterceptor(path)
	return nil
}

func SelectModule(namespace string, path string) core.Module {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(path)
	value, found := modules[key]
	if found {
		return value
	}
	return nil
}

func GetResourceSet(method string, level int) []core.Resource {
	value, found := resources[method][level]
	if found {
		return value
	}
	return nil
}

func AddResource(resource core.Resource) {
	method := resource.GetMethod()
	level := resource.GetLevel()
	rs := GetResourceSet(method, level)
	if rs == nil {
		rs = []core.Resource{}
	}
	resource.SetContext(&ResourceContext{
		namespace: resource.GetNamespace(),
	})
	resources[method][level] = append(rs, resource)
}

func RemoveResource(method string, path string) {
	if path != "" || strings.TrimSpace(path) != "" {
		path = meta.FormatPath(path)
		method = strings.ToUpper(method)
		level := len(strings.Split(path, "/")) - 1
		rs := GetResourceSet(method, level)
		if rs != nil && len(rs) > 0 {
			for i, v := range rs {
				if path == v.GetPath() {
					resources[method][level] = append(rs[:i], rs[i+1:]...)
					break
				}
			}
		}
	}
}

func SelectResource(method string, path string) core.Resource {
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	level := len(strings.Split(path, "/")) - 1
	rs := GetResourceSet(method, level)
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

func InitResources() error {
	log.Info("Start initializing resource information:")
	namespace := config.Namespace
	kvs, err := etcd.GetWithPrefix(meta.GetServiceRootPath(meta.FormatPath(namespace)))
	cnt := 0
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				log.Fatal(err)
			}
			if serv.Method == httpclient.MethodLocal {
				m := services.NewModule(serv)
				log.Infof("Service %s:%s:%s initialized successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), "LOCAL", meta.FormatPath(m.GetPath()))
				i := GenerateInterceptor(m)
				if i != nil {
					AddInterceptor(i)
					log.Infof(" --> Interceptor initialized successfully \n")
				} else {
					AddModule(m)
				}
			} else {
				res := NewResource(serv)
				AddResource(res)
				log.Infof("Service %s:%s:%s initialized successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), res.GetMethod(), meta.FormatPath(res.GetPath()))
			}
			cnt++
		}
		log.Infof("Resource initialization completed, a total of %d services were initialized. \n", cnt)
	}
	return err
}

func AddDataSource(ds *meta.DataSource) error {
	db, err := core.NewDB(ds)
	if err == nil {
		dbs[meta.FormatPath(db.GetNamespace())+meta.FormatPath(db.Name)] = db
	}
	return err
}

func RemoveDataSource(namespace string, name string) *core.RockDB {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(name)
	value, found := dbs[key]
	if found {
		delete(dbs, key)
		return value
	}
	return nil
}

func SelectDataSource(namespace string, name string) *core.RockDB {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(name)
	//value, found := dbs[key]
	//if found {
	//	return value
	//}
	return dbs[key]
}

func InitDbs() {
	log.Info("Start initializing datasource information:")
	namespace := config.Namespace
	cnt := 0
	ecnt := 0
	kvs, err := etcd.GetWithPrefix(meta.GetDataSourceRootPath(meta.FormatPath(namespace)))
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, kv := range kvs {
		ds, err := meta.NewDataSource(kv.Value)
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = AddDataSource(ds)
		namespace := cast.ToString(If(ds.Namespace == "", "default", ds.Namespace))
		if err != nil {
			log.Errorf("%s \n", errs.ErrorStyleFunc("DataSource "+namespace+":"+ds.Name+" initialized failed: "+err.Error()))
			ecnt++
		} else {
			log.Infof("DataSource %s:%s initialized successfully \n", namespace, ds.Name)
			cnt++
		}
	}
	log.Infof("DataSource initialization completed, a total of %d datasources were initialized, %d succeeded and %d failed. \n", cnt+ecnt, cnt, ecnt)
}

func metaWatcher(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
	splitted := strings.SplitN(string(key)[len(meta.GetMetaRootPath()):], "/", 5)
	namespace := splitted[1]
	metaType := splitted[2]
	if operType == etcd.CREATE {
		switch {
		case metaType == meta.ServiceDirectory[1:]:
			serv, err := meta.NewService(value)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if serv.Method == httpclient.MethodLocal {
				m := services.NewModule(serv)
				se := script.New()
				se.AddScript(script.GetSdk())
				se.AddScript("var exports={};")
				se.AddScript(m.GetSource())
				err := se.Run()
				if err != nil {
					log.Error(err)
					return
				}
				paths, err := se.GetMlVar("exports.interceptor.paths")
				if err != nil {
					log.Error(err)
					return
				}
				if paths != nil {
					if val, ok := paths.([]string); ok {
						level, _ := se.GetMlVar("exports.interceptor.paths")
						interceptor := services.NewInterceptor(m, val, cast.ToInt(level))
						requestHandler, err := se.GetMlFunc("exports.interceptor.requestHandler")
						if err != nil {
							log.Error(err)
							return
						}
						if requestHandler != nil {
							interceptor.SetRequestHandler(requestHandler)
						}
						responseHandler, err := se.GetMlFunc("exports.interceptor.responseHandler")
						if err != nil {
							log.Error(err)
							return
						}
						if responseHandler != nil {
							interceptor.SetResponseHandler(responseHandler)
						}
						AddInterceptor(interceptor)
						fmt.Printf("[INFO] Interceptor %s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "", 1), meta.FormatPath(m.GetPath()))
					}
				} else {
					AddModule(m)
					fmt.Printf("[INFO] Module %s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "", 1), meta.FormatPath(m.GetPath()))
				}
			} else {
				res := NewResource(serv)
				AddResource(res)
				fmt.Printf("[INFO] Service %s:%s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "", 1), res.GetMethod(), meta.FormatPath(res.GetPath()))
			}
			break
		case metaType == meta.DataSourceDirectory[1:]:
			ds, err := meta.NewDataSource(value)
			if err != nil {
				log.Fatalln(err)
			}
			err = AddDataSource(ds)
			if err != nil {
				fmt.Printf("[ERROR] %s \n", errs.ErrorStyleFunc("DataSource "+ds.Namespace+":"+ds.Name+" created failed: "+err.Error()))
			} else {
				fmt.Printf("[INFO] DataSource %s:%s created successfully \n", ds.Namespace, ds.Name)
			}
			break
		}
	} else if operType == etcd.MODIFY {
		switch {
		case metaType == meta.ServiceDirectory[1:]:
			serv, err := meta.NewService(value)
			if err != nil {
				log.Fatal(err)
				return
			}
			path := splitted[4]
			if serv.Method == httpclient.MethodLocal {
				m := services.NewModule(serv)
				log.Infof("Module %s: modified successfully \n", path)
				i := GenerateInterceptor(m)
				if i != nil {
					RemoveInterceptor(i.GetPath())
					AddInterceptor(i)
					log.Infof(" --> Interceptor modified successfully \n")
				} else {
					RemoveModule(serv.Namespace, serv.Path)
					AddModule(m)
				}
			} else {
				res := NewResource(serv)
				RemoveResource(serv.Method, serv.Path)
				AddResource(res)
				fmt.Printf("[INFO] Service %s:%s modified successfully \n", res.GetMethod(), path)
			}
			break
		}
	} else if operType == etcd.DELETE {
		switch {
		case metaType == meta.ServiceDirectory[1:]:
			method := splitted[3]
			path := splitted[4]
			if strings.ToUpper(method) == httpclient.MethodLocal {
				module := RemoveModule(namespace, path)
				if module != nil {
					fmt.Printf("[INFO] Module %s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), path)
				}
				i := RemoveInterceptor(path)
				if i != nil {
					log.Infof(" --> Interceptor deleted successfully \n")
				}
			} else {
				RemoveResource(method, path)
				fmt.Printf("[INFO] Service %s:%s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), strings.ToUpper(method), path)
			}
			break
		case metaType == meta.DataSourceDirectory[1:]:
			name := splitted[3]
			db := RemoveDataSource(namespace, name)
			if db != nil {
				fmt.Printf("[INFO] DataSource %s:%s deleted successfully \n", strings.Replace(db.GetNamespace(), meta.DefaultNamespace, "", 1), db.Name)
			}
			break
		}
	}
}

func GenerateInterceptor(module core.Module) *services.RockInterceptor {
	se := script.New()
	se.AddScript(script.GetSdk())
	se.AddScript(module.GetSource())
	err := se.Run()
	if err != nil {
		log.Error(err)
		return nil
	}
	paths, err := se.GetMlVar("exports.interceptor.paths")
	if err != nil {
		log.Error(err)
		return nil
	}
	if paths != nil {
		if val, ok := paths.([]string); ok {
			level, _ := se.GetMlVar("exports.interceptor.paths")
			interceptor := services.NewInterceptor(module.(*services.RockModule), val, cast.ToInt(level))
			requestHandler, err := se.GetMlFunc("exports.interceptor.requestHandler")
			if err != nil {
				log.Error(err)
				return nil
			}
			if requestHandler != nil {
				interceptor.SetRequestHandler(requestHandler)
			}
			responseHandler, err := se.GetMlFunc("exports.interceptor.responseHandler")
			if err != nil {
				log.Error(err)
				return nil
			}
			if responseHandler != nil {
				interceptor.SetResponseHandler(responseHandler)
			}
			return interceptor
		}
	}
	return nil
}

func WatchMeta() {
	etcd.WatchWithPrefix(meta.GetMetaRootPath(), metaWatcher)
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
