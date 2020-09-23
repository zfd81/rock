package server

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/zfd81/rock/errs"

	"github.com/zfd81/rock/script"

	"github.com/zfd81/rock/core"

	"github.com/pkg/errors"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/util/etcd"

	"github.com/zfd81/rock/http"
	"github.com/zfd81/rock/meta"
)

type Resource interface {
	SetContext(context core.Context)
	GetNamespace() string
	GetMethod() string
	GetPath() string
	GetRegexPath() string
	GetLevel() int
	GetPathParams() []*meta.Parameter
	AddPathParam(param *meta.Parameter)
	GetRequestParams() []*meta.Parameter
	AddRequestParam(param *meta.Parameter)
	Run() (log string, resp *http.Response, err error)
	Clear()
}

type ResourceContext struct {
	namespace string
}

func (c *ResourceContext) GetModule(path string) script.Module {
	return SelectModule(c.namespace, path)
}

func (c *ResourceContext) GetDataSource(name string) script.DB {
	return SelectDataSource(c.namespace, name)
}

func NewContext(namespace string) *ResourceContext {
	return &ResourceContext{
		namespace: namespace,
	}
}

var (
	config = conf.GetConfig()

	modules   = map[string]script.Module{} //模块映射
	resources = map[string]map[int][]Resource{
		http.MethodGet:    map[int][]Resource{}, // GET资源映射
		http.MethodPost:   map[int][]Resource{}, // POST资源映射
		http.MethodPut:    map[int][]Resource{}, // PUT资源映射
		http.MethodDelete: map[int][]Resource{}, // DELETE资源映射
	}

	dbs = map[string]*core.ParrotDB{}
)

func AddModule(module script.Module) {
	key := meta.FormatPath(module.GetNamespace()) + meta.FormatPath(module.GetPath())
	modules[key] = module
}

func RemoveModule(namespace string, path string) script.Module {
	if namespace == "" {
		namespace = meta.DefaultNamespace
	}
	key := meta.FormatPath(namespace) + meta.FormatPath(path)
	value, found := modules[key]
	if found {
		delete(modules, key)
		return value
	}
	return nil
}

func SelectModule(namespace string, path string) script.Module {
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

func GetResourceSet(method string, level int) []Resource {
	value, found := resources[method][level]
	if found {
		return value
	}
	return nil
}

func AddResource(resource Resource) {
	method := resource.GetMethod()
	level := resource.GetLevel()
	rs := GetResourceSet(method, level)
	if rs == nil {
		rs = []Resource{}
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

func SelectResource(method string, path string) Resource {
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
				for _, param := range resource.GetPathParams() {
					param.Value = pathFragments[param.Index]
				}
				resource.Clear()
				return resource
			}
		}
	}
	return nil
}

func InitResources() error {
	namespace := meta.DefaultNamespace
	if len(config.Namespaces) > 0 {
		namespace = config.Namespaces[0]
	}
	kvs, err := etcd.GetWithPrefix(meta.GetServiceRootPath() + meta.FormatPath(namespace))
	cnt := 0
	if err == nil {
		for _, kv := range kvs {
			serv, err := meta.NewService(kv.Value)
			if err != nil {
				log.Fatal(err)
			}
			if serv.Method == http.MethodLocal {
				m := core.NewModule(serv)
				AddModule(m)
			} else {
				res := core.NewResource(serv)
				AddResource(res)
				_, _, path := meta.ServicePath(string(kv.Key))
				fmt.Printf("[INFO] Service %s:%s initialized successfully \n", res.GetMethod(), path)
			}
			cnt++
		}
		fmt.Printf("[INFO] A total of %d services were initialized \n", cnt)
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

func RemoveDataSource(namespace string, name string) *core.ParrotDB {
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

func SelectDataSource(namespace string, name string) *core.ParrotDB {
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
	namespaces := config.Namespaces
	namespaces = append(namespaces, meta.DefaultNamespace)
	cnt := 0
	ecnt := 0
	for _, namespace := range namespaces {
		kvs, err := etcd.GetWithPrefix(meta.GetDataSourceRootPath() + meta.FormatPath(namespace))
		if err != nil {
			log.Fatalln(err.Error())
		}
		for _, kv := range kvs {
			ds, err := meta.NewDataSource(kv.Value)
			if err != nil {
				log.Fatalln(err.Error())
			}
			err = AddDataSource(ds)
			if err != nil {
				fmt.Printf("[ERROR] %s \n", errs.ErrorStyleFunc("DataSource "+ds.Namespace+":"+ds.Name+" initialized failed: "+err.Error()))
				ecnt++
			} else {
				fmt.Printf("[INFO] DataSource %s:%s initialized successfully \n", ds.Namespace, ds.Name)
				cnt++
			}
		}
	}
	fmt.Printf("[INFO] A total of %d datasources were initialized, %d succeeded and %d failed \n", cnt+ecnt, cnt, ecnt)
}

func metaWatcher(operType etcd.OperType, key []byte, value []byte, createRevision int64, modRevision int64, version int64) {
	full_path := meta.MetaPath(string(key))
	if operType == etcd.CREATE {
		switch {
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			serv, err := meta.NewService(value)
			if err != nil {
				log.Fatalln(err)
				return
			}
			namespace, _, path := meta.ServicePath(string(key))
			if serv.Method == http.MethodLocal {
				m := core.NewModule(serv)
				AddModule(m)
				fmt.Printf("[INFO] Module %s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), path)
			} else {
				res := core.NewResource(serv)
				AddResource(res)
				fmt.Printf("[INFO] Service %s:%s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), res.GetMethod(), path)
			}
			break
		case strings.HasPrefix(full_path, meta.DataSourceDirectory):
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
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			serv, err := meta.NewService(value)
			if err != nil {
				log.Fatal(err)
				return
			}
			if serv.Method == http.MethodLocal {
				m := core.NewModule(serv)
				RemoveModule(serv.Namespace, serv.Path)
				AddModule(m)
				_, _, path := meta.ServicePath(string(key))
				fmt.Printf("[INFO] Module %s: modified successfully \n", path)
			} else {
				res := core.NewResource(serv)
				RemoveResource(serv.Method, serv.Path)
				AddResource(res)
				_, _, path := meta.ServicePath(string(key))
				fmt.Printf("[INFO] Service %s:%s modified successfully \n", res.GetMethod(), path)
			}
			break
		}
	} else if operType == etcd.DELETE {
		switch {
		case strings.HasPrefix(full_path, meta.ServiceDirectory):
			namespace, method, path := meta.ServicePath(string(key))
			if strings.ToUpper(method) == http.MethodLocal {
				module := RemoveModule(namespace, path)
				if module != nil {
					fmt.Printf("[INFO] Module %s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), path)
				}
			} else {
				RemoveResource(method, path)
				fmt.Printf("[INFO] Service %s:%s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace, "", 1), strings.ToUpper(method), path)
			}
			break
		case strings.HasPrefix(full_path, meta.DataSourceDirectory):
			namespace, name := meta.DataSourcePath(string(key))
			db := RemoveDataSource(namespace, name)
			if db != nil {
				fmt.Printf("[INFO] DataSource %s:%s deleted successfully \n", strings.Replace(db.GetNamespace(), meta.DefaultNamespace, "", 1), db.Name)
			}
			break
		}
	}
}

func WatchMeta() {
	etcd.WatchWithPrefix(meta.GetMetaRootPath(), metaWatcher)
}
