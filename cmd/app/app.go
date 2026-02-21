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

	"github.com/gin-contrib/cors"
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

	// 4. 创建 Gin 引擎
	engine := gin.New()

	//engine.Use(cors.New(cors.Config{
	//	AllowOrigins: []string{
	//		"http://localhost:3004",
	//	},
	//	// 允许的 HTTP 方法
	//	AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	//	// 允许的请求头
	//	AllowHeaders: []string{
	//		"Origin",
	//		"Content-Type",
	//		"Accept",
	//		"Authorization",
	//		"X-Requested-With",
	//	},
	//	// 允许浏览器携带凭证（如 Cookie）
	//	AllowCredentials: true,
	//	// 预检请求的缓存时间（秒）
	//	MaxAge: 12 * time.Hour,
	//}))
	engine.Use(cors.Default()) // 允许所有源，但不支持 AllowCredentials:true

	// 5. 依赖注入
	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	recommendRepo := repo2.NewRecommendRepo(db)
	recommendService := service2.NewRecommendService(recommendRepo)
	recommendHandler := handler2.NewRecommendHandler(recommendService)

	// 6. 注册路由
	route.RegisterRoutes(engine, userHandler, recommendHandler) // 确认函数签名匹配

	// 返回 App 实例
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
