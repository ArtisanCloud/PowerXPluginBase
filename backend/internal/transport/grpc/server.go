package grpc

import (
	integrationtransport "github.com/ArtisanCloud/PowerXPlugin/internal/transport/grpc/integration"
	"google.golang.org/grpc"
)

// Registrar 聚合各模块的 gRPC 服务注册器。
type Registrar struct {
	Integration *integrationtransport.Server
}

// Register 将可用的 gRPC 服务注册到给定 server。
func Register(server *grpc.Server, registrar Registrar) {
	if server == nil {
		return
	}
	if registrar.Integration != nil {
		registrar.Integration.Register(server)
	}
}
