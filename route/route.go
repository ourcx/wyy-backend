// route/route.go
package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有路由，接收任意多个 Registrar
func RegisterRoutes(r *gin.Engine, registrars ...Registrar) {
	// 添加全局中间件（根据需要可放在这里或 main 中）
	r.Use(gin.Logger(), gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公共 API 路由组
	api := r.Group("/api")
	{
		// Ping 示例（也可以单独作为一个 Registrar）
		api.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// 让每个 Registrar 将自己的路由注册到 api 组下
		for _, reg := range registrars {
			reg.RegisterRoutes(api)
		}
	}
}
