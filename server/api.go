package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zfd81/parrot/core"
	"github.com/zfd81/parrot/meta/dai"

	"github.com/zfd81/parrot/server/env"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/rooster/types/container"
)

func param(c *gin.Context) container.Map {
	p := container.JsonMap{}
	err := c.ShouldBind(&p)
	if err != nil {
		return nil
	}
	return p
}

func Test(c *gin.Context) {
	serv := &meta.Service{}
	err := c.ShouldBind(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"log": err.Error(),
		})
		return
	}

	ins := env.New(serv)
	for _, param := range serv.Params {
		ins.SetParam(param.Name, param.Value)
	}

	log, err := ins.Run()
	if err != nil {
		log = log + fmt.Sprintf(env.LogFormat, time.Now().Format("2006-01-02 15:04:05.000")) + err.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"log":     log,
		"status":  ins.GetRespStatus(),
		"content": ins.GetRespContent(),
	})
}

func CreateService(c *gin.Context) {
	serv := &meta.Service{}
	err := c.ShouldBind(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	serv.Path = meta.FormatPath(serv.Path)
	err = dai.CreateService(serv)

	if err != nil {
		if err == core.ErrServExists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 101,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("Service %s created successfully", serv.Path),
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Method " + method + " not found",
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("Service %s deleted successfully", serv.Path),
	})
}

func ModifyService(c *gin.Context) {
	serv := &meta.Service{}
	err := c.ShouldBind(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}
	serv.Path = meta.FormatPath(serv.Path)
	err = dai.ModifyService(serv)

	if err != nil {
		if err == core.ErrServNotExist {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 102,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("Service %s modified successfully", serv.Path),
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Method " + method + " not found",
		})
		return
	}
	path := c.Param("path")
	serv, err := dai.GetService(namespace, m, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, serv)
}

func ListService(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	path := c.Param("path")
	servs, err := dai.ListService(namespace, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
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
	err = dai.CreateDataSource(ds)

	if err != nil {
		if err == core.ErrDsExists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 101,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("DataSource %s created successfully", ds.Name),
	})
}

func DeleteDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")
	err := dai.DeleteDataSource(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("DataSource %s deleted successfully", name),
	})
}

func ModifyDataSource(c *gin.Context) {
	ds := &meta.DataSource{}
	err := c.ShouldBind(ds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
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
		if err == core.ErrDsNotExist {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 102,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 100,
		"msg":  fmt.Sprintf("DataSource %s modified successfully", ds.Name),
	})
}

func FindDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")
	ds, err := dai.GetDataSource(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
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
