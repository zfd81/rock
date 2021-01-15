package server

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/server/services"
	"github.com/zfd81/rock/util/etcd"
)

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
				m := services.NewModule(serv).SetEnvironment(env)
				i := m.GenerateInterceptor()
				log.Infof("Module %s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), meta.FormatPath(m.GetPath()))
				if i != nil {
					AddInterceptor(i)
					log.Infof(" --> Interceptor created successfully \n")
				} else {
					env.AddModule(m)
				}
			} else {
				res := services.NewResource(serv).SetEnvironment(env)
				env.AddResource(res)
				log.Infof("Service %s:%s:%s created successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), res.GetMethod(), meta.FormatPath(res.GetPath()))
			}
			break
		case metaType == meta.DataSourceDirectory[1:]:
			ds, err := meta.NewDataSource(value)
			if err != nil {
				log.Fatalln(err)
			}
			err = env.AddDataSource(ds)
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
				m := services.NewModule(serv).SetEnvironment(env)
				log.Infof("Module %s: modified successfully \n", path)
				i := m.GenerateInterceptor()
				if i != nil {
					RemoveInterceptor(i.GetPath())
					AddInterceptor(i)
					log.Infof(" --> Interceptor modified successfully \n")
				} else {
					env.RemoveModule(serv.Namespace, serv.Path)
					env.AddModule(m)
				}
			} else {
				res := services.NewResource(serv).SetEnvironment(env)
				env.RemoveResource(serv.Method, serv.Path)
				env.AddResource(res)
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
				module := env.RemoveModule(namespace, path)
				if module != nil {
					log.Infof("Module %s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), path)
				}
				i := RemoveInterceptor(path)
				if i != nil {
					log.Infof(" --> Interceptor deleted successfully \n")
				}
			} else {
				env.RemoveResource(method, path)
				log.Infof("Service %s:%s:%s deleted successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), strings.ToUpper(method), path)
			}
			break
		case metaType == meta.DataSourceDirectory[1:]:
			name := splitted[3]
			db := env.RemoveDataSource(namespace, name)
			if db != nil {
				log.Infof("DataSource %s:%s deleted successfully \n", strings.Replace(db.GetNamespace(), meta.DefaultNamespace[1:], "default", 1), db.GetName())
			}
			break
		}
	}
}

func WatchMeta() {
	etcd.WatchWithPrefix(meta.GetMetaRootPath(), metaWatcher)
}
