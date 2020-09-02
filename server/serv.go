package server

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/zfd81/parrot/core"

	"github.com/spf13/cast"

	"github.com/zfd81/parrot/meta"

	"github.com/pkg/errors"

	"github.com/zfd81/parrot/server/env"

	"github.com/gin-gonic/gin"
)

func CallGetService(c *gin.Context) {
	path := c.Param("path")
	resource := selectResource(path, env.GetResources())
	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		if err == core.ErrParamBad {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Service request parameter error",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
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
	resource := selectResource(path, env.PostResources())

	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		if err == core.ErrParamBad {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Service request parameter error",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
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
	c.JSON(http.StatusOK, resp.Content)
}

func CallPutService(c *gin.Context) {
	path := c.Param("path")
	resource := selectResource(path, env.PutResources())

	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		if err == core.ErrParamBad {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Service request parameter error",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resource)
}

func CallDeleteService(c *gin.Context) {
	path := c.Param("path")
	resource := selectResource(path, env.DeleteResources())

	if resource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		if err == core.ErrParamBad {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Service request parameter error",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resource)
}

func selectResource(path string, resourceMap map[int][]env.Resource) env.Resource {
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	level := len(strings.Split(path, "/")) - 1
	resources := resourceMap[level]
	if resources != nil {
		for _, resource := range resources {
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

func wrapParam(c *gin.Context, resource env.Resource) error {
	if len(resource.GetRequestParams()) > 0 {
		p := param(c)
		if p == nil || p.Empty() {
			return core.ErrParamBad
		}
		for _, param := range resource.GetRequestParams() {
			val, found := p.Get(param.Name)
			if !found {
				return errors.New("parameter " + param.Name + " not found")
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
	parrot := e.Group("/")
	{
		parrot.GET("/*path", CallGetService)
		parrot.POST("/*path", CallPostService)
		parrot.PUT("/*path", CallPutService)
		parrot.DELETE("/*path", CallDeleteService)
	}
	return e
}
