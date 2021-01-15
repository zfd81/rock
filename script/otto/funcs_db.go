package otto

import (
	"reflect"
	"strings"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rock/core"
)

func DBQuery(env core.Context) func(datasource string, query string, arg interface{}, pageNumber int, pageSize int) []container.Map {
	return func(datasource string, query string, arg interface{}, pageNumber int, pageSize int) []container.Map {
		db := env.GetDataSource(datasource) //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			throwException("Data source[%s] not found", datasource)
		}
		sql := strings.TrimSpace(query) //获取SQL
		if sql == "" {
			throwException("SQL statement cannot be empty")
		}
		if pageNumber > 0 {
			if pageSize < 1 {
				pageSize = 10
			}
			l, err := db.QueryMapList(sql, arg, pageNumber, pageSize)
			if err != nil {
				throwException(err.Error())
			}
			return l
		} else {
			r, err := db.Query(sql, arg)
			if err != nil {
				throwException(err.Error())
			}
			l, err := r.MapListScan()
			if err != nil {
				throwException(err.Error())
			}
			return l
		}
	}
}

func DBQueryOne(env core.Context) func(datasource string, query string, arg interface{}) container.Map {
	return func(datasource string, query string, arg interface{}) container.Map {
		db := env.GetDataSource(datasource) //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			throwException("Data source[%s] not found", datasource)
		}
		sql := strings.TrimSpace(query) //获取SQL
		if sql == "" {
			throwException("SQL statement cannot be empty")
		}
		m, err := db.QueryMap(sql, arg)
		if err != nil {
			throwException(err.Error())
		}
		return m
	}
}

func DBSave(env core.Context) func(datasource string, table string, arg interface{}) int64 {
	return func(datasource string, table string, arg interface{}) int64 {
		db := env.GetDataSource(datasource) //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			throwException("Data source[%s] not found", datasource)
		}
		table = strings.TrimSpace(table)
		if table == "" {
			throwException("Table name cannot be empty")
		}
		if arg == nil {
			throwException("Parameter cannot be empty")
		}
		m, ok := arg.(map[string]interface{})
		if ok {
			num, err := db.Save(m, table)
			if err != nil {
				throwException(err.Error())
			}
			return num
		} else {
			l, ok := arg.([]interface{})
			if ok {
				num, err := db.BatchSave(l, table)
				if err != nil {
					throwException(err.Error())
				}
				return num
			} else {
				l, ok := arg.([]map[string]interface{})
				if ok {
					num, err := db.BatchSave(SliceParam(l), table)
					if err != nil {
						throwException(err.Error())
					}
					return num
				} else {
					throwException("Parameter data type error")
					return -1
				}
			}
		}
	}
}

func DBExec(env core.Context) func(datasource string, query string, arg interface{}) int64 {
	return func(datasource string, query string, arg interface{}) int64 {
		db := env.GetDataSource(datasource) //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			throwException("Data source[%s] not found", datasource)
		}
		sql := strings.TrimSpace(query) //获取SQL
		if sql == "" {
			throwException("SQL statement cannot be empty")
		}
		//v, ok := arg.([]interface{})
		//if !ok {
		//	v, ok := arg.([]map[string]interface{})
		//	if ok {
		//		arg = v
		//	}
		//} else {
		//	arr := []interface{}{}
		//	for _, i := range v {
		//		arr = append(arr, i.(map[string]interface{}))
		//	}
		//	arg = arr
		//}
		num, err := db.Exec(sql, arg)
		if err != nil {
			throwException(err.Error())
		}
		return num
	}
}

func SliceParam(args []map[string]interface{}) []interface{} {
	param := make([]interface{}, len(args))
	for i, v := range args {
		param[i] = v
	}
	return param
}
