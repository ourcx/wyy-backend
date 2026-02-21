package handler

import (
	service "wyy/internal/service/discover"
	"wyy/utils"

	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	RecommendService *service.RecommendService
}

func NewRecommendHandler(recommendService *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{RecommendService: recommendService}
}

// RegisterRoutes 实现 route.Registrar 接口
func (h *RecommendHandler) RegisterRoutes(r gin.IRouter) {
	recommends := r.Group("/recommends")
	{
		recommends.GET("/banners", h.getReBanners)
	}
}

// getReBanners 返回 Banner 假数据
// @Summary      获取推荐 Banner
// @Description  返回 Banner 图片 URL 列表（当前为假数据）
// @Tags         推荐模块
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "成功返回 banners 列表"
// @Router       /api/recommends/banners [get]
// @Example      { "banners": ["https://p5.music.126.net/obj/wonDlsKUwrLClGjCm8Kx/78326360541/d092/1425/a80f/5ee0ef793063a2f70d3001d5cacc517c.jpg"] }
func (h *RecommendHandler) getReBanners(c *gin.Context) {
	banners := []string{
		"https://p5.music.126.net/obj/wonDlsKUwrLClGjCm8Kx/78326360541/d092/1425/a80f/5ee0ef793063a2f70d3001d5cacc517c.jpg",
	}
	//包装的返回值对象
	utils.Success(c, banners)
}
