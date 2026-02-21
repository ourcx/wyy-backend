// route/interface.go
package route

import "github.com/gin-gonic/gin"

// Registrar 定义路由注册器接口
type Registrar interface {
	// RegisterRoutes 将路由注册到给定的路由组上
	RegisterRoutes(r gin.IRouter)
}
