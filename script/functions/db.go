package functions

import (
	"reflect"
	"strings"

	"github.com/zfd81/parrot/errs"

	"github.com/zfd81/parrot/script"

	"github.com/robertkrimen/otto"
)

func DBQuery(env script.Environment) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) (value otto.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return ErrorResult(call, "Data source["+name+"] not found")
		}
		sql_v := call.Argument(1)
		if !sql_v.IsString() {
			return ErrorResult(call, "SQL statement cannot be empty")
		}
		sql := strings.TrimSpace(sql_v.String()) //获取SQL
		var arg interface{}
		pageNumber := -1 //当前页码
		pageSize := 10   //页面数据量
		arg_v := call.Argument(2)
		if arg_v.IsObject() {
			arg_v, err := arg_v.Export()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsString() {
			arg_v, err := arg_v.ToString()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsNumber() {
			arg_int, err := arg_v.ToInteger()
			if err == nil {
				arg = arg_int
			} else {
				arg_float, err := arg_v.ToFloat()
				if err != nil {
					return ErrorResult(call, err.Error())
				}
				arg = arg_float
			}
		}
		pageNumber_v := call.Argument(3)
		if pageNumber_v.IsDefined() {
			if !pageNumber_v.IsNumber() {
				return ErrorResult(call, "Parameter pageNumber data type error")
			}
			pageNumber_v, err := pageNumber_v.ToInteger()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			pageNumber = int(pageNumber_v)
			pageSize_v := call.Argument(4)
			if pageSize_v.IsNumber() {
				pageSize_v, err := pageSize_v.ToInteger()
				if err != nil {
					return ErrorResult(call, err.Error())
				}
				pageSize = int(pageSize_v)
			}
			l, err := db.QueryMapList(sql, arg, pageNumber, pageSize)
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			return Result(call, l)
		} else {
			r, err := db.Query(sql, arg)
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			l, err := r.MapListScan()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			return Result(call, l)
		}
		return
	}
}

func DBQueryOne(env script.Environment) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) (value otto.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return ErrorResult(call, "Data source["+name+"] not found")
		}
		sql_v := call.Argument(1)
		if !sql_v.IsString() {
			return ErrorResult(call, "SQL statement cannot be empty")
		}
		sql := strings.TrimSpace(sql_v.String()) //获取SQL
		var arg interface{}
		arg_v := call.Argument(2)
		if arg_v.IsObject() {
			arg_v, err := arg_v.Export()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsString() {
			arg_v, err := arg_v.ToString()
			if err != nil {
				return ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsNumber() {
			arg_int, err := arg_v.ToInteger()
			if err == nil {
				arg = arg_int
			} else {
				arg_float, err := arg_v.ToFloat()
				if err != nil {
					return ErrorResult(call, err.Error())
				}
				arg = arg_float
			}
		}
		m, err := db.QueryMap(sql, arg)
		if err != nil {
			return ErrorResult(call, err.Error())
		}
		return Result(call, m)
	}
}

func Result(call otto.FunctionCall, data interface{}) (value otto.Value) {
	result := &script.Result{
		Code: 200,
		Data: data,
	}
	value, _ = call.Otto.ToValue(result)
	return
}

func ErrorResult(call otto.FunctionCall, err string) (value otto.Value) {
	result := &script.Result{
		Code:    400,
		Message: errs.ErrorStyleFunc(err),
	}
	value, _ = call.Otto.ToValue(result)
	return
}
