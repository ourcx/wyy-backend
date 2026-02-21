package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wyy/route"

	"wyy/internal/config"
	"wyy/internal/handler"
	handler2 "wyy/internal/handler/discover"
	"wyy/internal/repo"
	repo2 "wyy/internal/repo/discover"
	"wyy/internal/service"
	service2 "wyy/internal/service/discover"

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

	// 3. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 4. 创建 Gin 引擎并添加中间件
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	// 5. 依赖注入
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	recommendRepo := repo2.NewRecommendRepo(db)
	recommendService := service2.NewRecommendService(recommendRepo)
	recommendHandler := handler2.NewRecommendHandler(recommendService)

	// 6. 注册路由（以分组方式）
	route.RegisterRoutes(engine, userHandler, recommendHandler)

	return &App{
		cfg:    cfg,
		db:     db,
		router: engine,
	}, nil
}

func (a *App) Close() error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}
