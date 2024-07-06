package usecase

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/internal/projects"
)

type tasksUC struct {
	tasksRepo projects.TasksRepository
}

func NewTasksUseCase(tasksRepo projects.TasksRepository) projects.TasksUseCase {
	return tasksUC{
		tasksRepo: tasksRepo,
	}
}

func (t tasksUC) Get(ctx context.Context, projectID int64) ([]*models.Task, error) {
	return t.tasksRepo.Get(ctx, projectID)
}

func (t tasksUC) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	return t.tasksRepo.Create(ctx, task)
}

func (t tasksUC) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	return t.tasksRepo.Update(ctx, task)
}

func (t tasksUC) Delete(ctx context.Context, taskID int64) error {
	return t.tasksRepo.Delete(ctx, taskID)
}

func (t tasksUC) Start(ctx context.Context, taskID, userID int64) error {
	return t.tasksRepo.Start(ctx, taskID, userID)
}

func (t tasksUC) Stop(ctx context.Context, taskID, userID int64) error {
	return t.tasksRepo.Stop(ctx, taskID, userID)
}

func (t tasksUC) GetMembers(ctx context.Context, taskID int64) ([]*models.User, error) {
	return t.tasksRepo.GetMembers(ctx, taskID)
}

func (t tasksUC) AddMember(ctx context.Context, taskID, userID int64) error {
	return t.tasksRepo.AddMember(ctx, taskID, userID)
}

func (t tasksUC) DeleteMember(ctx context.Context, taskID, userID int64) error {
	return t.tasksRepo.DeleteMember(ctx, taskID, userID)
}

func (t tasksUC) IsMember(ctx context.Context, taskID, userID int64) (bool, error) {
	return t.tasksRepo.IsMember(ctx, taskID, userID)
}
