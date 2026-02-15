package main

import (
	"fmt"
	"wyy/internal/config"
	//"wyy/internal/handler"
	"wyy/internal/repo"
	//"wyy/internal/service"
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

	//// 3. 初始化路由（同时完成依赖注入）
	//router := setupRouter(db)

	return &App{
		cfg: cfg,
		db:  db,
		//router: router,
	}, nil
}

//func setupRouter(db *gorm.DB) *gin.Engine {
//	// 初始化各层
//	userRepo := repo.NewUserRepo(db)
//	userService := service.NewUserService(userRepo)
//	userHandler := handler.NewUserHandler(userService)
//
//	router := gin.Default()
//	router.POST("/register", userHandler.Register)
//	// 其他路由...
//	return router
//}

func (a *App) Run() error {
	return a.router.Run(fmt.Sprintf(":%d", a.cfg.Server.Port))
}
