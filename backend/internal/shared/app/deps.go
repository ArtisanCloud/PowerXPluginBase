package app

import (
	"context"
	"github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"
	"gorm.io/gorm"
)

type Deps struct {
	DB           *gorm.DB
	Ctx          *context.Context
	PowerXClient *client.PowerXServiceClient
}
