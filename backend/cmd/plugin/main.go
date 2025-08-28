package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scrum-plugin/internal/config"
	"scrum-plugin/internal/db"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/router"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger.Init(cfg.LogLevel)
	logger.Info("Starting PowerX Scrum Plugin...")

	// 连接数据库
	if err := db.Connect(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	// 注意：不再在应用启动时执行任何迁移操作
	// 数据库迁移应该通过独立的 cmd/database/migrate 命令执行

	// 设置路由
	r := router.New(cfg)
	engine := r.Setup()

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: engine,

		// 超时配置
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,

		// 最大头部大小
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		logger.WithField("addr", cfg.BindAddr).Info("Starting HTTP server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed")
	}
}
