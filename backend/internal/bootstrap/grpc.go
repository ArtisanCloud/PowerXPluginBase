package bootstrap

import (
	"context"
	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
)

func BootstrapGRPCClient(ctx context.Context, cfg *config.GRPCUpstream) *client.PowerXServiceClient {

	pxc, err := client.NewPowerXServiceClient(ctx, cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize PowerX gRPC client")
	}

	logger.Info("PowerX gRPC client initialized successfully")

	return pxc
}
