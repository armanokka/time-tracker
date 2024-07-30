package repository

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTasksRepository_Create(t *testing.T) {
	tasksRepo, db, mock, err := newMockTasksRepo()
	require.NoError(t, err)
	defer db.Close()

	task := getTestTask()
	mock.ExpectQuery(createTaskQuery).WithArgs(task.Name, task.Description, task.ProjectID).WillReturnRows(
		sqlmock.NewRows(task.Columns()).AddRow(task.Fields()...),
	)

	gotTask, err := tasksRepo.Create(context.Background(), task)
	assert.Nil(t, err)
	assert.Equal(t, task, gotTask)
}

func TestTasksRepository_Delete(t *testing.T) {
	tasksRepo, db, mock, err := newMockTasksRepo()
	require.NoError(t, err)
	defer db.Close()

	task := getTestTask()

	mock.ExpectExec(deleteTaskQuery).WithArgs(task.ID).WillReturnResult(sqlmock.NewResult(4, 3))
	assert.Nil(t, tasksRepo.Delete(context.Background(), task.ID))

	mock.ExpectExec(deleteTaskQuery).WithArgs(task.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	assert.NotNil(t, tasksRepo.Delete(context.Background(), task.ID))
}

func TestTasksRepository_AddMember(t *testing.T) {
	tasksRepo, db, mock, err := newMockTasksRepo()
	require.NoError(t, err)
	defer db.Close()

	task := getTestTask()
	var userID int64 = 9

	mock.ExpectExec(addTaskMemberQuery).WithArgs(task.ID, userID).WillReturnResult(sqlmock.NewResult(1, 1))
	assert.Nil(t, tasksRepo.AddMember(context.Background(), task.ID, userID))

	mock.ExpectExec(addTaskMemberQuery).WithArgs(task.ID, userID).WillReturnResult(sqlmock.NewResult(1, 0))
	assert.NotNil(t, tasksRepo.AddMember(context.Background(), task.ID, userID))

	mock.ExpectExec(addTaskMemberQuery).WithArgs(task.ID, userID).WillReturnError(fmt.Errorf("some error"))
	assert.NotNil(t, tasksRepo.AddMember(context.Background(), task.ID, userID))
}

func TestTasksRepository_Get(t *testing.T) {
	tasksRepo, db, mock, err := newMockTasksRepo()
	require.NoError(t, err)
	defer db.Close()

	task := getTestTask()
	var projectID int64 = 1

	tasks := []*models.Task{task}

	mock.ExpectQuery(getTotalTasks).WithArgs(projectID).WillReturnRows(
		sqlmock.NewRows([]string{"count"}).AddRow(len(tasks)),
	)
	mock.ExpectQuery(selectTasks).WithArgs(projectID).WillReturnRows(
		sqlmock.NewRows(task.Columns()).AddRow(task.Fields()...),
	)

	gotTasks, err := tasksRepo.Get(context.Background(), projectID)
	assert.Nil(t, err)
	assert.Equal(t, tasks, gotTasks)
}
