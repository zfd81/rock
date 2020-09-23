package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zfd81/rooster/types/container"
)

type ApiResponse struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"msg,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func param(c *gin.Context) (container.Map, error) {
	p := container.JsonMap{}
	err := c.ShouldBind(&p)
	return p, err
}

func Router() http.Handler {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())
	api := e.Group("/rock")
	{
		api.POST("/test/analysis", TestAnalysis)                //服务分析
		api.POST("/test", Test)                                 //测试服务
		api.POST("/serv", CreateService)                        //创建服务
		api.DELETE("/serv/method/:method/*path", DeleteService) //删除服务
		api.PUT("/serv", ModifyService)                         //修改服务
		api.GET("/serv/method/:method/*path", FindService)      //查询单个服务
		api.GET("/serv/list/*path", ListService)                //查询服务列表

		api.POST("/ds", CreateDataSource)              //创建数据源
		api.DELETE("/ds/name/:name", DeleteDataSource) //删除数据源
		api.PUT("/ds", ModifyDataSource)               //修改数据源
		api.GET("/ds/name/:name", FindDataSource)      //查询单个数据源
		api.GET("/ds/list", ListDataSource)            //查询数据源列表
	}
	return e
}
