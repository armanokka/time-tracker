package projects

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
)

type RedisRepository interface {
	GetProject(ctx context.Context, projectID int64) (*models.Project, error)
	SetProject(ctx context.Context, project *models.Project, seconds int64) error
	DeleteProject(ctx context.Context, projectID int64) error
}
