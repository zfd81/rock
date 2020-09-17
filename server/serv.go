package server

import (
	"net/http"
	"strings"

	"github.com/zfd81/parrot/conf"

	"github.com/zfd81/parrot/errs"

	"github.com/spf13/cast"

	"github.com/zfd81/parrot/meta"

	"github.com/gin-gonic/gin"
)

func CallGetService(c *gin.Context) {
	path := c.Param("path")
	resource := SelectResource(http.MethodGet, path)
	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log, resp, err := resource.Run()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  log,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, resp.Data)
}

func CallPostService(c *gin.Context) {
	path := c.Param("path")
	resource := SelectResource(http.MethodPost, path)

	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log, resp, err := resource.Run()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  log,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, resp.Data)
}

func CallPutService(c *gin.Context) {
	path := c.Param("path")
	resource := SelectResource(http.MethodPut, path)

	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log, resp, err := resource.Run()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  log,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, resp.Data)
}

func CallDeleteService(c *gin.Context) {
	path := c.Param("path")
	resource := SelectResource(http.MethodDelete, path)

	if resource == nil {
		c.JSON(http.StatusNotFound, errs.New(404, "Target service["+path+"] not exist."))
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log, resp, err := resource.Run()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  log,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, resp.Data)
}

func wrapParam(c *gin.Context, resource Resource) error {
	if len(resource.GetRequestParams()) > 0 {
		p, err := param(c)
		if err != nil {
			return err
		}
		if p == nil || p.Empty() {
			return errs.New(errs.ErrParamBad)
		}
		for _, param := range resource.GetRequestParams() {
			val, found := p.Get(param.Name)
			if !found {
				return errs.New(errs.ErrParamNotFound, param.Name)
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
	return nil
}

func ParrotRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())
	parrot := e.Group(conf.GetConfig().ServiceName)
	{
		parrot.GET("/*path", CallGetService)
		parrot.POST("/*path", CallPostService)
		parrot.PUT("/*path", CallPutService)
		parrot.DELETE("/*path", CallDeleteService)
	}
	return e
}
