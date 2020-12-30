package server

import (
	"net/http"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/errs"

	"github.com/gin-gonic/gin"
)

func param(c *gin.Context) (container.Map, error) {
	p := container.JsonMap{}
	err := c.ShouldBind(&p)
	return p, err
}

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
			return errs.NewError(err)
		}
		if p == nil || p.Empty() {
			return errs.New(errs.ErrParamBad)
		}
		for _, param := range resource.GetRequestParams() {
			val, found := p.Get(param.Name)
			if !found {
				return errs.New(errs.ErrParamNotFound, param.Name)
			}
			if err = param.SetValue(val); err != nil {
				return errs.New(errs.ErrParamBad, err.Error())
			}
		}
	}
	return nil
}

func Router() http.Handler {
	e := gin.New()
	e.Use(Logger())
	parrot := e.Group(conf.GetConfig().ServiceName)
	{
		parrot.GET("/*path", CallGetService)
		parrot.POST("/*path", CallPostService)
		parrot.PUT("/*path", CallPutService)
		parrot.DELETE("/*path", CallDeleteService)
	}
	return e
}
