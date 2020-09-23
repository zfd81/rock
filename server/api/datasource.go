package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/rock/errs"
	"github.com/zfd81/rock/meta"
	"github.com/zfd81/rock/meta/dai"
)

func CreateDataSource(c *gin.Context) {
	ds := &meta.DataSource{}
	err := c.ShouldBind(ds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	//数据源内容监测
	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "[DataSource information cannot be empty]"))
		return
	}
	if err = dai.CreateDataSource(ds); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("DataSource %s created successfully", ds.Name),
	})
}

func DeleteDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")                        //从Path中获得数据源名称
	ds := &meta.DataSource{
		Namespace: namespace,
		Name:      name,
	}
	if err := dai.DeleteDataSource(ds); err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("DataSource %s deleted successfully", name),
	})
}

func ModifyDataSource(c *gin.Context) {
	ds := &meta.DataSource{}
	err := c.ShouldBind(ds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	//数据源内容监测
	if ds.Name == "" ||
		ds.Driver == "" ||
		ds.Host == "" ||
		ds.Port < 100 ||
		ds.User == "" ||
		ds.Password == "" ||
		ds.Database == "" {
		c.JSON(http.StatusBadRequest, errs.New(errs.ErrParamBad, "[DataSource information cannot be empty]"))
		return
	}
	if err = dai.ModifyDataSource(ds); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("DataSource %s modified successfully", ds.Name),
	})
}

func FindDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	name := c.Param("name")                        //从Path中获得数据源名称
	ds, err := dai.GetDataSource(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError(err))
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Data:       ds,
	})
}

func ListDataSource(c *gin.Context) {
	namespace := c.Request.Header.Get("namespace") //从Header中获得命名空间
	dses, err := dai.ListDataSource(namespace, "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		StatusCode: 200,
		Data:       dses,
	})
}
