package server

import (
	"sort"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/server/services"
	"github.com/zfd81/rock/util/etcd"
)

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
		err = env.AddDataSource(ds)
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

func InitResources() error {
	namespace := config.Namespace
	kvs, err := etcd.GetWithPrefix(meta.GetServiceRootPath(meta.FormatPath(namespace)))
	if err != nil {
		return err
	}
	var kvslice KeyValueSlice = kvs
	sort.Sort(kvslice)
	err = InitModules(kvslice)
	if err != nil {
		return err
	}
	return InitServices(kvs)
}

func InitModules(kvs []*mvccpb.KeyValue) error {
	log.Info("Start initializing modules information:")
	namespace := config.Namespace
	cnt := 0
	for _, kv := range kvs {
		serv, err := meta.NewService(kv.Value)
		if err != nil {
			log.Fatal(err)
		}
		if serv.Method == httpclient.MethodLocal {
			m := services.NewModule(serv).SetEnvironment(env)
			i := m.GenerateInterceptor()
			log.Infof("Module %s:%s initialized successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), meta.FormatPath(m.GetPath()))
			if i != nil {
				AddInterceptor(i)
				log.Infof(">>>>>> Interceptor initialized successfully \n")
			} else {
				env.AddModule(m)
			}
			cnt++
		}
	}
	log.Infof("Modules initialization completed, a total of %d modules were initialized. \n", cnt)
	return nil
}

func InitServices(kvs []*mvccpb.KeyValue) error {
	log.Info("Start initializing services information:")
	namespace := config.Namespace
	cnt := 0
	for _, kv := range kvs {
		serv, err := meta.NewService(kv.Value)
		if err != nil {
			log.Fatal(err)
		}
		if serv.Method != httpclient.MethodLocal {
			res := services.NewResource(serv).SetEnvironment(env)
			env.AddResource(res)
			log.Infof("Service %s:%s:%s initialized successfully \n", strings.Replace(namespace, meta.DefaultNamespace[1:], "default", 1), res.GetMethod(), meta.FormatPath(res.GetPath()))
			cnt++
		}
	}
	log.Infof("Services initialization completed, a total of %d services were initialized. \n", cnt)
	return nil
}

type KeyValueSlice []*mvccpb.KeyValue

func (kvs KeyValueSlice) Len() int {
	return len(kvs)
}

func (kvs KeyValueSlice) Less(i, j int) bool {
	return kvs[i].CreateRevision < kvs[j].CreateRevision
}

func (kvs KeyValueSlice) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}
