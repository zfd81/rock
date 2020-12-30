package server

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Log struct {
	// Start time
	StartTime time.Time

	// EndTime shows the time after the server returns a response.
	EndTime time.Time

	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration

	RequestURI string

	// 传入服务器请求的协议版本
	Proto string

	// 用户代理
	UserAgent string

	// ClientIP equals Context's ClientIP method.
	ClientIP string

	// Method is the HTTP method given to the request.
	Method string

	// Path is a path the client requests.
	Path string

	// StatusCode is HTTP response code.
	StatusCode int

	// BodySize is the size of the Response Body
	BodySize int

	//用户标识
	userId string
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		l := &Log{}
		end := time.Now()
		l.StartTime = start
		l.EndTime = end
		l.Latency = end.Sub(start)
		l.RequestURI = c.Request.RequestURI
		l.Proto = c.Request.Proto
		l.UserAgent = c.Request.UserAgent()
		l.ClientIP = c.ClientIP()
		l.Method = c.Request.Method
		l.StatusCode = c.Writer.Status()
		l.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		l.Path = path
		l.userId = "abcd"
		log.Infof("status=%d elapsed=%v client=%s method=%s path=%s", l.StatusCode, l.Latency, l.ClientIP, l.Method, l.Path)
	}
}
