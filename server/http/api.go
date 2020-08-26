package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zfd81/parrot/meta/dai"

	"github.com/zfd81/parrot/server/env"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/rooster/types/container"
)

func param(c *gin.Context) container.Map {
	p := container.JsonMap{}
	c.ShouldBind(&p)
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
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}
	serv.Name = meta.FormatServiceName(serv.Name)
	cnt, err := dai.CreateService(serv)

	if err != nil {
		if cnt == -1 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 1234,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1234,
		"msg":  "ok",
	})
}

func DeleteService(c *gin.Context) {
	name := c.Param("name")
	err := dai.DeleteService(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1234,
		"msg":  "ok",
	})
}

func ModifyService(c *gin.Context) {
	serv := &meta.Service{}
	err := c.ShouldBind(serv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}
	err = dai.ModifyService(serv)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1234,
		"msg":  "ok",
	})
}

func FindService(c *gin.Context) {
	name := c.Param("name")
	serv, err := dai.GetService(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, serv)
}

func ListService(c *gin.Context) {
	name := c.Param("name")
	if name == "*" {
		name = ""
	}
	servs, err := dai.ListService(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1234,
			"msg":  err.Error(),
		})
		return
	}

	names := make([]string, 0, 50)
	for _, serv := range servs {
		names = append(names, serv.Name)
	}

	c.JSON(http.StatusOK, names)
}
