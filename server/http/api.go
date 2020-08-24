package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/parrot/meta"
	"github.com/zfd81/parrot/script"
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
			"msg": err.Error(),
		})
		return
	}
	se := script.New()
	for _, param := range serv.Params {
		se.AddVar(param.Name, param.Value)
	}
	se.AddScript(serv.Script)
	log, err := se.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"log": log,
	})
}
