package usecase

import (
	"context"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tasksUC struct {
	tasksRepo projects.TasksRepository
	tracer    trace.Tracer
}

func NewTasksUseCase(tasksRepo projects.TasksRepository) projects.TasksUseCase {
	return tasksUC{
		tasksRepo: tasksRepo,
		tracer:    otel.GetTracerProvider().Tracer("api"),
	}
}

func (t tasksUC) Get(ctx context.Context, taskID int64) ([]*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Get")
	defer span.End()

	return t.tasksRepo.Get(ctx, taskID)
}

func (t tasksUC) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Create")
	defer span.End()

	return t.tasksRepo.Create(ctx, task)
}

func (t tasksUC) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Update")
	defer span.End()

	return t.tasksRepo.Update(ctx, task)
}

func (t tasksUC) Delete(ctx context.Context, taskID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Delete")
	defer span.End()

	return t.tasksRepo.Delete(ctx, taskID)
}

func (t tasksUC) Start(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Start")
	defer span.End()

	return t.tasksRepo.Start(ctx, taskID, userID)
}

func (t tasksUC) Stop(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.Stop")
	defer span.End()

	return t.tasksRepo.Stop(ctx, taskID, userID)
}

func (t tasksUC) GetMembers(ctx context.Context, taskID int64) ([]*models.User, error) {
	ctx, span := t.tracer.Start(ctx, "tasksUC.GetMembers")
	defer span.End()

	return t.tasksRepo.GetMembers(ctx, taskID)
}

func (t tasksUC) AddMember(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.AddMember")
	defer span.End()

	return t.tasksRepo.AddMember(ctx, taskID, userID)
}

func (t tasksUC) DeleteMember(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.DeleteMember")
	defer span.End()

	return t.tasksRepo.DeleteMember(ctx, taskID, userID)
}

func (t tasksUC) IsMember(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksUC.IsMember")
	defer span.End()

	return t.tasksRepo.IsMember(ctx, taskID, userID)
}
