package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"scrum-plugin/internal/config"
	"scrum-plugin/internal/db"
	"scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/grpc/server"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/router"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
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

	// 初始化 PowerX gRPC 客户端
	ctx := context.Background()
	pxc, err := client.NewPowerX(ctx, cfg.GRPCUpstream)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize PowerX gRPC client")
	}
	defer func() {
		if err := pxc.Close(); err != nil {
			logger.WithError(err).Error("Failed to close PowerX gRPC client")
		}
	}()

	logger.Info("PowerX gRPC client initialized successfully")

	// 设置路由
	r := router.New(cfg)
	r.SetPowerXClient(pxc) // 设置 PowerX 客户端
	engine := r.Setup()

	// 创建 gRPC 服务器（可选）
	gs, err := server.New(ctx, cfg.GRPCServer)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create gRPC server")
	}

	// 创建 HTTP 服务器
	httpServer := &http.Server{
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

	// 使用 errgroup 并发启动服务器
	g, ctx := errgroup.WithContext(context.Background())

	// 启动 HTTP 服务器
	g.Go(func() error {
		logger.WithField("addr", cfg.BindAddr).Info("Starting HTTP server...")
		return httpServer.ListenAndServe()
	})

	// 启动 gRPC 服务器（如果启用）
	if gs != nil {
		g.Go(func() error {
			return gs.Serve(ctx)
		})
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的 goroutine 中等待信号
	go func() {
		<-quit
		logger.Info("Shutting down servers...")

		// 优雅关闭 HTTP 服务器
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.WithError(err).Error("HTTP server forced to shutdown")
		} else {
			logger.Info("HTTP server shutdown completed")
		}

		// gRPC 服务器会通过 context 取消自动关闭
	}()

	// 等待服务器启动失败或优雅关闭
	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.WithError(err).Error("Server error")
		os.Exit(1)
	}

	logger.Info("All servers shutdown completed")
}
