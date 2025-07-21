package main

import (
	"AUV/config"
	"AUV/db"
	"AUV/middleware"
	"AUV/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	if err := db.InitDB(cfg); err != nil {
		log.Fatal("Failed to init DB:", err)
	}

	// 自动迁移
	db.AutoMigrate()

	// 初始化Gin
	r := gin.Default()

	// 配置静态文件路由
	r.Static(cfg.Server.StaticPrefix, "./public")

	// 放行所有跨域请求
	r.Use(middleware.Cors())

	// ip限流
	r.Use(middleware.RateLimit())

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
