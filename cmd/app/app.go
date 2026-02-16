package main

import (
	"fmt"
	"log"
	"wyy/index"
	"wyy/internal/config"
	"wyy/internal/handler"
	"wyy/internal/repo"
	"wyy/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	router *gin.Engine
}

func NewApp(cfgPath string) (*App, error) {
	// 1. 加载配置
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// 2. 初始化数据库
	db, err := repo.NewDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	// 3. 设置 Gin 模式（从配置读取）
	gin.SetMode(cfg.Server.Mode) // 例如 "debug" 或 "release"

	// 4. 创建 Gin 引擎（不使用默认中间件，我们手动添加）
	engine := gin.New()

	// 5. 依赖注入：创建各层实例
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 6. 注册路由（将 handler 传入）
	route.RegisterRoutes(engine, userHandler)

	return &App{
		cfg:    cfg,
		db:     db,
		router: engine,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	return a.router.Run(addr)
}
