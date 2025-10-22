package console

import "github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"

// Service aggregates shared dependencies for admin console workflows.
type Service struct {
	deps *app.Deps
}

// NewService creates a new console Service placeholder.
func NewService(deps *app.Deps) *Service {
	return &Service{deps: deps}
}
