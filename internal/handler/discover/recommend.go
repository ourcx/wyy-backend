package handler

import (
	service "wyy/internal/service/discover"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	RecommendService *service.RecommendService
}

func NewRecommendHandler(userService *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{RecommendService: userService}
}

func (h *RecommendHandler) RegisterRoutes(r gin.IRouter) {
	//recommends := r.Group("/recommends")
	{
		//recommends.GET("/", h.GetRecommendations)
		// ...
	}
}

// 从数据库拿到 banner 的图片链接
func (h *RecommendHandler) getReBanners(c *gin.Context) {

}
