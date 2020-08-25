package http

import (
	"fmt"
	"net/http"
	"time"

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
