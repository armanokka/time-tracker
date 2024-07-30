package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tasksRepository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewTasksRepository(db *sqlx.DB) projects.TasksRepository {
	return tasksRepository{db: db, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (t tasksRepository) Get(ctx context.Context, projectID int64) ([]*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Get")
	defer span.End()

	var totalCount int
	if err := t.db.GetContext(ctx, &totalCount, getTotalTasks, projectID); err != nil {
		return nil, err
	}

	rows, err := t.db.QueryxContext(ctx, selectTasks, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.Task, 0, totalCount)
	for rows.Next() {
		var task models.Task
		if err = rows.StructScan(&task); err != nil {
			return nil, err
		}
		results = append(results, &task)
	}
	return results, nil
}

func (t tasksRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Get")
	defer span.End()

	return task, t.db.QueryRowxContext(ctx, createTaskQuery, task.Name, task.Description,
		task.ProjectID).StructScan(task)
}

func (t tasksRepository) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Update")
	defer span.End()

	return task, t.db.QueryRowxContext(ctx, updateTaskQuery, task.Name, task.Description,
		task.ID).StructScan(task)
}

func (t tasksRepository) Delete(ctx context.Context, taskID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Delete")
	defer span.End()

	result, err := t.db.ExecContext(ctx, deleteTaskQuery, taskID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (t tasksRepository) Start(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Start")
	defer span.End()

	var count int
	if err := t.db.GetContext(ctx, &count, getActiveUserTasksQuery, userID, taskID); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("task already started")
	}

	result, err := t.db.ExecContext(ctx, startTaskQuery, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (t tasksRepository) Stop(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.Stop")
	defer span.End()

	result, err := t.db.ExecContext(ctx, endTaskQuery, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (t tasksRepository) GetMembers(ctx context.Context, taskID int64) ([]*models.User, error) {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.GetMembers")
	defer span.End()

	var totalCount int
	if err := t.db.GetContext(ctx, &totalCount, getTotalTaskMembersQuery, taskID); err != nil {
		return nil, err
	}

	rows, err := t.db.QueryxContext(ctx, getTaskMembersQuery, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0, totalCount)
	for rows.Next() {
		user := &models.User{}
		if err = rows.StructScan(user); err != nil {
			return nil, err
		}
		user.Sanitize()
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (t tasksRepository) AddMember(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.AddMember")
	defer span.End()

	result, err := t.db.ExecContext(ctx, addTaskMemberQuery, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (t tasksRepository) DeleteMember(ctx context.Context, taskID, userID int64) error {
	ctx, span := t.tracer.Start(ctx, "tasksRepository.AddMember")
	defer span.End()

	result, err := t.db.ExecContext(ctx, deleteTaskMemberQuery, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (t tasksRepository) IsMember(ctx context.Context, taskID, userID int64) error {
	result, err := t.db.ExecContext(ctx, isTaskMemberQuery, taskID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
