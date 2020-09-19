package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zfd81/rock/core"

	"github.com/spf13/cast"

	"github.com/zfd81/rock/script"

	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta/dai"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rooster/types/container"
)

func param(c *gin.Context) (container.Map, error) {
	p := container.JsonMap{}
	err := c.ShouldBind(&p)
	return p, err
}

func TestAnalysis(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	definition, code := SplitSource(p.GetString("source"))
	serv := &meta.Service{}
	err = wrapService(definition, serv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	serv.Script = code
	c.JSON(http.StatusOK, serv)
}

func Test(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	definition, code := SplitSource(p.GetString("source"))
	serv := &meta.Service{}
	err = wrapService(definition, serv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	serv.Script = code
	res := core.NewResource(serv)
	ps, found := p.Get("params")
	if found {
		params := cast.ToStringMap(ps)
		for _, param := range res.GetPathParams() {
			param.Value = cast.ToString(params[param.Name])
		}
		for _, param := range res.GetRequestParams() {
			val, found := params[param.Name]
			if !found {
				c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamNotFound, param.Name))
				return
			}
			if strings.ToUpper(param.DataType) == meta.DataTypeString {
				param.Value = cast.ToString(val)
			} else if strings.ToUpper(param.DataType) == meta.DataTypeInteger {
				param.Value = cast.ToInt(val)
			} else if strings.ToUpper(param.DataType) == meta.DataTypeBool {
				param.Value = cast.ToBool(val)
			} else if strings.ToUpper(param.DataType) == meta.DataTypeMap {
				param.Value = cast.ToStringMap(val)
			} else if strings.ToUpper(param.DataType) == meta.DataTypeArray {
				param.Value = cast.ToStringSlice(val)
			}
		}
	}
	res.SetContext(&ResourceContext{
		namespace: res.GetNamespace(),
	})
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
	definition, code := SplitSource(p.GetString("source"))
	serv := &meta.Service{}
	err = wrapService(definition, serv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Name = p.GetString("name")
	serv.Script = code
	err = dai.CreateService(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Service %s created successfully", serv.Path),
	})
}

func DeleteService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	method := c.Param("method")
	m := strings.ToUpper(method)
	if m != http.MethodGet &&
		m != http.MethodPost &&
		m != http.MethodPut &&
		m != http.MethodDelete {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "Method "+method+" not found"))
		return
	}
	path := c.Param("path")
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
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Service %s deleted successfully", serv.Path),
	})
}

func ModifyService(c *gin.Context) {
	p, err := param(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	definition, code := SplitSource(p.GetString("source"))
	serv := &meta.Service{}
	err = wrapService(definition, serv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	serv.Script = code
	err = dai.ModifyService(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Service %s modified successfully", serv.Path),
	})
}

func FindService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	method := c.Param("method")
	m := strings.ToUpper(method)
	if m != http.MethodGet &&
		m != http.MethodPost &&
		m != http.MethodPut &&
		m != http.MethodDelete {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "Method "+method+" not found"))
		return
	}
	path := c.Param("path")
	serv, err := dai.GetService(namespace, m, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, serv)
}

func ListService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	path := c.Param("path")
	servs, err := dai.ListService(namespace, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}

	paths := make([]string, 0, 50)
	for _, serv := range servs {
		paths = append(paths, strings.ToUpper(serv.Method)+":"+serv.Path)
	}

	c.JSON(http.StatusOK, paths)
}

func CreateDataSource(c *gin.Context) {
	ds := &meta.DataSource{}
	err := c.ShouldBind(ds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}

	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "DataSource information cannot be empty"))
		return
	}
	err = dai.CreateDataSource(ds)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("DataSource %s created successfully", ds.Name),
	})
}

func DeleteDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")
	err := dai.DeleteDataSource(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("DataSource %s deleted successfully", name),
	})
}

func ModifyDataSource(c *gin.Context) {
	ds := &meta.DataSource{}
	err := c.ShouldBind(ds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "DataSource information cannot be empty",
		})
		return
	}
	err = dai.ModifyDataSource(ds)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("DataSource %s modified successfully", ds.Name),
	})
}

func FindDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")
	ds, err := dai.GetDataSource(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, ds)
}

func ListDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	dses, err := dai.ListDataSource(namespace, "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dses)
}

func ApiRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())
	api := e.Group("/parrot")
	{
		api.POST("/test", Test)
		api.POST("/test/analysis", TestAnalysis)
		api.POST("/serv", CreateService)
		api.DELETE("/serv/method/:method/*path", DeleteService)
		api.PUT("/serv", ModifyService)
		api.GET("/serv/method/:method/*path", FindService)
		api.GET("/serv/list/*path", ListService)

		api.POST("/ds", CreateDataSource)
		api.DELETE("/ds/name/:name", DeleteDataSource)
		api.PUT("/ds", ModifyDataSource)
		api.GET("/ds/name/:name", FindDataSource)
		api.GET("/ds/list", ListDataSource)
	}
	return e
}

func SplitSource(source string) (string, string) {
	start := strings.Index(source, "$.define(")
	if start == -1 {
		return "", source
	}
	end := strings.Index(source[start:], "})")
	if end == -1 {
		return "", source
	}
	return source[start : end+3], source[end+3:]
}

func wrapService(definition string, serv *meta.Service) error {
	if definition == "" {
		return errs.New(errs.ErrParamBad, "Missing service definition")
	}
	se := script.New()
	se.SetScript(definition)
	se.Run()
	data, err := se.GetVar("__serv_definition")
	if err != nil {
		return errs.New(errs.ErrParamBad, "Service definition error")
	}
	val, ok := data.(map[string]interface{})
	if !ok {
		return errs.New(errs.ErrParamBad, "Service definition error")
	}
	namespace := cast.ToString(val["namespace"])
	path := cast.ToString(val["path"])
	if path == "" {
		return errs.New(errs.ErrParamBad, "Service path not found")
	}
	method := cast.ToString(val["method"])
	if method == "" {
		return errs.New(errs.ErrParamBad, "Service method not found")
	}
	m := strings.ToUpper(method)
	if m != http.MethodGet &&
		m != http.MethodPost &&
		m != http.MethodPut &&
		m != http.MethodDelete {
		return errs.New(errs.ErrParamBad, "Service method["+method+"] error")
	}
	params := val["params"]
	serv.Namespace = namespace
	serv.Path = path
	serv.Method = method
	if params != nil {
		ps, ok := params.([]map[string]interface{})
		if !ok {
			return errs.New(errs.ErrParamBad, "Service parameter definition error")
		}
		for _, param := range ps {
			serv.AddParam(cast.ToString(param["name"]), cast.ToString(param["dataType"]))
		}
	}
	return nil
}
