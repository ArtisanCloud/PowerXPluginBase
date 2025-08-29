package services

import (
	"context"
	"fmt"

	"scrum-plugin/internal/domain/models"
	"scrum-plugin/internal/domain/repository"
	"scrum-plugin/internal/logger"

	"gorm.io/gorm"
)

// SprintService Sprint 服务结构体
type SprintService struct {
	sprintRepo *repository.SprintRepository
}

// NewSprintService 创建 Sprint 服务实例
func NewSprintService(db *gorm.DB) *SprintService {
	sprintRepo := repository.NewSprintRepository(db)
	return &SprintService{
		sprintRepo: sprintRepo,
	}
}

// GetSprint 根据ID获取Sprint
func (s *SprintService) GetSprint(ctx context.Context, tenantID int64, id uint64) (*models.Sprint, error) {
	log := logger.ServiceLogger("sprint").WithContext(ctx)

	sprint, err := s.sprintRepo.GetById(ctx, id, nil)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"sprint_id": id,
		}).Error("Failed to get sprint by ID")
		return nil, fmt.Errorf("failed to get sprint: %w", err)
	}

	if sprint == nil {
		log.WithFields(logger.Fields{
			"sprint_id": id,
		}).Warn("Sprint not found")
		return nil, fmt.Errorf("sprint not found")
	}

	return sprint, nil
}

// ListSprints 获取Sprint列表
func (s *SprintService) ListSprints(ctx context.Context, tenantID int64, opts *repository.SprintListOptions) ([]*models.Sprint, int64, error) {
	log := logger.ServiceLogger("sprint").WithContext(ctx)

	sprints, total, err := s.sprintRepo.List(ctx, opts)
	if err != nil {
		log.WithError(err).Error("Failed to list sprints")
		return nil, 0, fmt.Errorf("failed to list sprints: %w", err)
	}

	log.WithFields(logger.Fields{
		"count": len(sprints),
		"total": total,
	}).Debug("Sprints listed successfully")

	return sprints, total, nil
}
