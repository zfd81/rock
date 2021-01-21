package server

import (
	"net/http"

	"github.com/zfd81/rock/httpclient"
	"github.com/zfd81/rock/server/services"
	"github.com/zfd81/rooster/types/container"

	"github.com/gin-gonic/gin"
)

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		chain := GetInterceptorChain()
		if chain.Len() > 0 {
			normal := true
			req := httpclient.NewRequest(c.Request)
			resp := httpclient.NewResponse()
			resp.StatusCode = http.StatusOK
			s := container.NewArrayStack()
			path := req.GetPath()
			for _, i := range chain {
				if i.Matches(path) {
					s.Push(i)
					ok, err := i.Request(req, resp)
					if err != nil {
						c.JSON(http.StatusBadRequest, err)
						return
					}
					if !ok {
						normal = false
						break
					}
				}
			}
			header := resp.Header
			for k, _ := range header {
				c.Header(k, header.Get(k))
			}
			if normal {
				c.Next()
			}
			for !s.Empty() {
				s.Peek()
				i, _ := s.Pop()
				ok, err := i.(*services.RockInterceptor).Response(req, resp)
				if err != nil {
					c.JSON(http.StatusBadRequest, err)
					return
				}
				if !ok {
					break
				}
			}
			if resp.Data != nil {
				c.JSON(resp.StatusCode, resp.Data)
			}
			c.Abort()
		} else {
			c.Next()
		}
	}
}
