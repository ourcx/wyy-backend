package route

import (
	"net/http"
	"wyy/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
// 这里接收需要的 handler 实例，实现依赖注入
func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	// 全局中间件（可根据配置决定是否开启）
	r.Use(gin.Logger())   // 请求日志
	r.Use(gin.Recovery()) // panic 恢复

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 示例：用户相关路由
	users := r.Group("/api/users")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
		users.GET("/:id", userHandler.GetUser)
	}

	// 你可以继续添加其他模块的路由分组
}
