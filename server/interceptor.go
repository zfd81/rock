package server

import (
	"net/http"

	"github.com/zfd81/rock/server/services"
	"github.com/zfd81/rooster/types/container"

	"github.com/gin-gonic/gin"
)

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		chain := GetInterceptorChain()
		if chain.Len() > 0 {
			normal := true
			resp := &http.Response{}
			s := container.NewArrayStack()
			path := c.Request.URL.Path
			for _, i := range chain {
				if i.Matches(path) {
					s.Push(i)
					ok, err := i.Request(c.Request, resp)
					if err != nil {
						return
					}
					if !ok {
						normal = false
						break
					}
				}
			}
			if normal {
				c.Next()
			}
			for !s.Empty() {
				s.Peek()
				i, _ := s.Pop()
				i.(*services.RockInterceptor).Response(c.Request, resp)
			}
		} else {
			c.Next()
		}
	}
}
