package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zfd81/rock/server"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zfd81/rock/core"
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/meta/dai"
	"github.com/zfd81/rock/script"
)

func TestAnalysis(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	source := p.GetString("source")
	serv, err := SourceAnalysis(source)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	c.JSON(http.StatusOK, serv)
}

func Test(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	source := p.GetString("source")
	serv, err := SourceAnalysis(source)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	serv.Source = source
	if serv.Method == "LOCAL" {

	}
	res := core.NewResource(serv)
	ps, found := p.Get("params")
	if found {
		params := cast.ToStringMap(ps)
		for _, param := range res.GetPathParams() {
			param.SetValue(cast.ToString(params[param.Name]))
		}
		for _, param := range res.GetRequestParams() {
			val, found := params[param.Name]
			if !found {
				c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamNotFound, param.Name))
				return
			}
			err = param.SetValue(val)
			if err != nil {
				c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "Parameter data type error"))
				return
			}
		}
	}
	res.SetContext(server.NewContext(res.GetNamespace()))
	log, resp, err := res.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"log": log,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, gin.H{
		"log":    log,
		"header": resp.Header,
		"data":   resp.Data,
	})
}

func CreateService(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	serv, err := SourceAnalysis(p.GetString("source"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	serv.Source = p.GetString("source")
	err = dai.CreateService(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Service %s created successfully", serv.Path),
	})
}

func DeleteService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	method := c.Param("method")                    //从Path中获得方法
	m := strings.ToUpper(method)
	if m != http.MethodGet && m != http.MethodPost &&
		m != http.MethodPut && m != http.MethodDelete &&
		m != "LOCAL" {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "Method "+method+" not found"))
		return
	}
	path := c.Param("path") //从Path中获得服务的访问路径
	serv := &meta.Service{
		Namespace: namespace,
		Method:    method,
		Path:      path,
	}
	err := dai.DeleteService(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Service %s deleted successfully", serv.Path),
	})
}

func ModifyService(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	serv, err := SourceAnalysis(p.GetString("source"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Source = p.GetString("source")
	err = dai.ModifyService(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Service %s modified successfully", serv.Path),
	})
}

func FindService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	method := c.Param("method")                    //从Path中获得方法
	m := strings.ToUpper(method)
	if m != http.MethodGet && m != http.MethodPost &&
		m != http.MethodPut && m != http.MethodDelete &&
		m != "LOCAL" {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "Method "+method+" not found"))
		return
	}
	path := c.Param("path") //从Path中获得服务的访问路径
	serv, err := dai.GetService(namespace, m, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	if serv != nil {
		serv.Source = ""
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Data:       serv,
	})
}

func ListService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	path := c.Param("path")
	servs, err := dai.ListService(namespace, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	for _, serv := range servs {
		serv.Source = ""
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Data:       servs,
	})
}

func SourceAnalysis(source string) (*meta.Service, error) {
	var definition string
	start := strings.Index(source, "$.define(")
	if start != -1 {
		end := strings.Index(source[start:], "})")
		if end == -1 {
			return nil, errs.New(errs.ErrParamBad, "Service definition error")
		}
		definition = source[start : end+3]
	}
	var namespace string
	var path string
	var method string
	serv := &meta.Service{}
	se := script.New()
	if definition != "" {
		se.SetScript(definition)
		err := se.Run()
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Service definition error:"+err.Error())
		}
		data, err := se.GetVar("__serv_definition")
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Service definition error:"+err.Error())
		}
		val, ok := data.(map[string]interface{})
		if !ok {
			return nil, errs.New(errs.ErrParamBad, "Service definition error")
		}
		namespace = cast.ToString(val["namespace"])
		path = cast.ToString(val["path"])
		if path == "" {
			return nil, errs.New(errs.ErrParamBad, "Service path not found")
		}
		method = cast.ToString(val["method"])
		if method == "" {
			return nil, errs.New(errs.ErrParamBad, "Service method not found")
		}
		m := strings.ToUpper(method)
		if m != http.MethodGet && m != http.MethodPost &&
			m != http.MethodPut && m != http.MethodDelete {
			return nil, errs.New(errs.ErrParamBad, "Service method["+method+"] error")
		}
		params := val["params"]
		if params != nil {
			ps, ok := params.([]map[string]interface{})
			if !ok {
				return nil, errs.New(errs.ErrParamBad, "Service parameters definition error")
			}
			for _, param := range ps {
				serv.AddParam(cast.ToString(param["name"]), cast.ToString(param["dataType"]))
			}
		}
	} else {
		se.AddScript(se.GetSdk())
		se.AddScript("var module={};")
		se.AddScript(source)
		err := se.Run()
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Module definition error:"+err.Error())
		}
		value, err := se.GetVar("module")
		if err != nil {
			return nil, errs.New(errs.ErrParamBad, "Module definition error:"+err.Error())
		}
		module, ok := value.(map[string]interface{})
		if !ok {
			return nil, errs.New(errs.ErrParamBad, "Module definition error")
		}
		value = module["exports"]
		exports, ok := value.(map[string]interface{})
		if !ok {
			return nil, errs.New(errs.ErrParamBad, "Module definition error")
		}
		namespace = cast.ToString(exports["namespace"])
		path = cast.ToString(exports["path"])
		if path == "" {
			return nil, errs.New(errs.ErrParamBad, "Module path not found")
		}
		method = "LOCAL"
	}
	serv.Namespace = namespace
	serv.Path = path
	serv.Method = method
	return serv, nil
}
