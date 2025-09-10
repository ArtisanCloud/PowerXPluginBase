package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"

	cfgpkg "github.com/ArtisanCloud/PowerXPlugin/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server 插件 gRPC 服务器
type Server struct {
	*grpc.Server
	lis    net.Listener
	config *cfgpkg.GRPCServer
}

// New 创建新的插件 gRPC 服务器
func NewGRPCServer(ctx context.Context, c *cfgpkg.GRPCServer) (*Server, error) {
	if !c.Enable {
		logger.Info("gRPC server is disabled")
		return nil, nil
	}

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", c.Addr, err)
	}

	var opts []grpc.ServerOption

	if c.UseTLS {
		if c.Cert == "" || c.Key == "" {
			return nil, fmt.Errorf("TLS is enabled but cert or key is missing")
		}
		creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(creds))
		logger.Info("gRPC server TLS enabled")
	} else {
		// 明确声明：开发期不加 TLS
		logger.Warn("gRPC server running without TLS (development mode)")
		_ = insecure.NewCredentials()
	}

	s := grpc.NewServer(opts...)

	// 注册健康检查服务
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)

	// 设置服务健康状态
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("note-plugin", healthpb.HealthCheckResponse_SERVING)

	// 注册反射服务（开发和调试用）
	reflection.Register(s)

	// TODO: 在这里注册你的插件 gRPC 服务
	// 例如：pluginv1.RegisterNotePluginServiceServer(s, NewNoteServer(deps))

	//logger.WithField("address", lis.Addr().String()).Info("gRPC server configured")
	logger.Info("gRPC server configured")

	return &Server{
		Server: s,
		lis:    lis,
		config: c,
	}, nil
}

// Serve 启动 gRPC 服务器
func (s *Server) Serve(ctx context.Context) error {
	logger.WithField("address", s.lis.Addr().String()).Info("Starting gRPC server")

	// 在单独的 goroutine 中监听上下文取消
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	return s.Server.Serve(s.lis)
}

// GetListenAddr 获取监听地址
func (s *Server) GetListenAddr() string {
	if s.lis != nil {
		return s.lis.Addr().String()
	}
	return s.config.Addr
}

// IsServing 检查服务器是否在运行
func (s *Server) IsServing() bool {
	return s.lis != nil
}

// TODO: 当定义插件自己的 proto 服务时，在这里实现服务逻辑
// 例如：
//
// type NoteServer struct {
// 	pluginv1.UnimplementedNotePluginServiceServer
// 	noteService *services.NoteService
// 	// 其他依赖
// }
//
// func NewNoteServer(deps *SomeDependencies) *NoteServer {
// 	return &NoteServer{
// 		noteService: deps.NoteService,
// 	}
// }
//
// func (s *NoteServer) CreateNote(ctx context.Context, req *pluginv1.CreateNoteRequest) (*pluginv1.CreateNoteResponse, error) {
// 	// 实现创建任务逻辑
// 	return &pluginv1.CreateNoteResponse{}, nil
// }
//
// func (s *NoteServer) GetNote(ctx context.Context, req *pluginv1.GetNoteRequest) (*pluginv1.GetNoteResponse, error) {
// 	// 实现获取任务逻辑
// 	return &pluginv1.GetNoteResponse{}, nil
// }
