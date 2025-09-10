package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerXPlugin/internal/bootstrap"
	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/grpc/server"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	"github.com/ArtisanCloud/PowerXPlugin/internal/router"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	ctx := context.Background()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化插件
	queryDB, err := bootstrap.BootstrapPlugin(ctx, cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to bootstrap plugin")
	}

	// 初始化 PowerX gRPC Client 客户端
	pxc := bootstrap.BootstrapGRPCClient(ctx, cfg.GRPCUpstream)

	deps := &app.Deps{
		DB:           queryDB,
		Ctx:          &ctx,
		PowerXClient: pxc,
	}

	// 设置 gin engine 路由
	r := router.NewRouter(cfg, deps)
	engine := r.Setup()

	// 创建 gRPC 服务器（可选）
	gs, err := server.NewGRPCServer(ctx, cfg.GRPCServer)
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
