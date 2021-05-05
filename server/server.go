package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/zfd81/rock/core"

	"github.com/zfd81/rooster/types/container"

	"github.com/zfd81/rock/conf"

	"github.com/zfd81/rock/errs"

	"github.com/gin-gonic/gin"
)

func param(c *gin.Context) (container.Map, error) {
	p := container.JsonMap{}
	queryMap := c.Request.URL.Query() //获得URL中的参数
	for k, vals := range queryMap {
		if len(vals) == 1 {
			p.Put(k, vals[0])
		} else {
			p.Put(k, vals)
		}
	}
	err := c.ShouldBind(&p)
	return p, err
}

func CallGetService(c *gin.Context) {
	path := c.Param("path")
	resource := env.SelectResource(http.MethodGet, path)
	if resource == nil {
		log.Error("Target service[" + path + "] not exist.")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	loginfo, resp, err := resource.Run()

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  loginfo,
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
	resource := env.SelectResource(http.MethodPost, path)

	if resource == nil {
		log.Error("Target service[" + path + "] not exist.")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	loginfo, resp, err := resource.Run()

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  loginfo,
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
	resource := env.SelectResource(http.MethodPut, path)

	if resource == nil {
		log.Error("Target service[" + path + "] not exist.")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Target service[" + path + "] not exist.",
		})
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	loginfo, resp, err := resource.Run()

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  loginfo,
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
	resource := env.SelectResource(http.MethodDelete, path)

	if resource == nil {
		log.Error("Target service[" + path + "] not exist.")
		c.JSON(http.StatusNotFound, errs.New(404, "Target service["+path+"] not exist."))
		return
	}

	err := wrapParam(c, resource)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	loginfo, resp, err := resource.Run()

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  loginfo,
		})
		return
	}
	for k, v := range resp.Header {
		c.Header(k, v)
	}
	c.JSON(http.StatusOK, resp.Data)
}

func wrapParam(c *gin.Context, resource core.Resource) error {
	if len(resource.GetParams()) > 0 {
		p, err := param(c)
		if err != nil {
			return errs.NewError(err)
		}
		//if p == nil || p.Empty() {
		//	return errs.New(errs.ErrParamBad)
		//}
		for _, param := range resource.GetParams() {
			if param.IsRequestScope() {
				val, found := p.Get(param.Name)
				if !found {
					return errs.New(errs.ErrParamNotFound, param.Name)
				}
				if err = param.SetValue(val); err != nil {
					return errs.New(errs.ErrParamBad, err.Error())
				}
			} else if param.IsHeaderScope() {
				val := c.Request.Header.Get(param.Name)
				if val == "" {
					return errs.New(errs.ErrParamNotFound, param.Name)
				}
				if err = param.SetValue(val); err != nil {
					return errs.New(errs.ErrParamBad, err.Error())
				}
			}
		}
	}
	return nil
}

func Router() http.Handler {
	e := gin.New()
	e.Use(Logger(), Interceptor())
	rock := e.Group(conf.GetConfig().ServiceName)
	{
		rock.GET("/*path", CallGetService)
		rock.POST("/*path", CallPostService)
		rock.PUT("/*path", CallPutService)
		rock.DELETE("/*path", CallDeleteService)
	}
	return e
}
