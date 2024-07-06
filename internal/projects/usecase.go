package projects

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, project *models.Project) (*models.Project, error)
	GetByID(ctx context.Context, projectID int64) (*models.Project, error)
	Update(ctx context.Context, updates *models.Project) (*models.Project, error)
	Delete(ctx context.Context, projectID int64) error

	IsOwner(ctx context.Context, projectID, userID int64) (bool, error)
	IsMember(ctx context.Context, projectID, userID int64) (bool, error)

	GetMembers(ctx context.Context, projectID int64) ([]*models.User, error)
	AddMember(ctx context.Context, projectID, userID int64) error
	RemoveMember(ctx context.Context, projectID, userID int64) error
	GetMemberProductivity(ctx context.Context, projectID, userID int64) ([]models.UserProductivity, error)
}

type TasksUseCase interface {
	Get(ctx context.Context, projectID int64) ([]*models.Task, error)
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) (*models.Task, error)
	Delete(ctx context.Context, taskID int64) error

	Start(ctx context.Context, taskID, userID int64) error
	Stop(ctx context.Context, taskID, userID int64) error

	GetMembers(ctx context.Context, taskID int64) ([]*models.User, error)
	AddMember(ctx context.Context, taskID, userID int64) error
	DeleteMember(ctx context.Context, taskID, userID int64) error
	IsMember(ctx context.Context, taskID, userID int64) (bool, error)
}
