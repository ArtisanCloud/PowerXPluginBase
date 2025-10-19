package integration

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Server 占位 integration gRPC 服务。
type Server struct {
	logger *logrus.Entry
}

// NewServer 构造 integration gRPC 服务。
func NewServer(logger *logrus.Entry) *Server {
	return &Server{logger: logger}
}

// Register 将 integration 服务注册到 gRPC Server。
func (s *Server) Register(registrar grpc.ServiceRegistrar) {
	if registrar == nil {
		return
	}
	// TODO: 注册实际的 integration proto 服务。
	s.log().Debug("integration gRPC server placeholder registered")
}

func (s *Server) log() *logrus.Entry {
	if s.logger != nil {
		return s.logger
	}
	return logrus.WithField("component", "integration.grpc")
}
