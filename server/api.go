package server

import (
	"fmt"
	"net/http"

	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta/dai"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/rock/meta"
)

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
	ds := &meta.DataSource{
		Namespace: namespace,
		Name:      name,
	}
	err := dai.DeleteDataSource(ds)
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
	dses, err := dai.ListDataSource(namespace, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 999,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dses)
}
