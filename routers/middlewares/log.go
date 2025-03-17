package middlewares

import (
	"time"

	"github.com/turbo-uid/hots/globals"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		globals.GoLogger.Infof("completed handling request: path=%s method=%s httpcode=%d %s", c.Request.URL.Path, c.Request.Method, c.Writer.Status(), time.Since(start).String())
	}
}
