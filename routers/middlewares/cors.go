package middlewares

import (
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func CorsReq() gin.HandlerFunc {
	config := cors.Config{
		// 允许所有域名进行跨域请求
		AllowOrigins: []string{"*"},
		// 允许所有HTTP方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		// 允许所有自定义头
		AllowHeaders: []string{"*"},
		// 暴露给客户端的自定义头
		ExposeHeaders: []string{"*"},
		// 是否允许发送Cookie
		AllowCredentials: true,
		// 预检请求的最大缓存时间, 0表示不缓存
		MaxAge: 0,
	}
	return cors.New(config)
}
