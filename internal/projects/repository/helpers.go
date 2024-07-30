package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

func getTestProject() *models.Project {
	description := "lorem ipsum"
	return &models.Project{
		ID:          15,
		Name:        "Some project",
		Description: &description,
		CreatorID:   10,
	}
}

func getTestTask() *models.Task {
	return &models.Task{
		ID:          1,
		Name:        "Lorem",
		Description: "Ipsum doromet",
		ProjectID:   4,
		Finished:    false,
	}
}

// SetupRedis launches local Redis instance via testcontainers. Returned testcontainers.Container MUST be terminated
func SetupRedis(ctx context.Context) (testcontainers.Container, *redis.Client) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	return redisC, redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
}

func newMockProjectsRepo() (projects.Repository, *sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, nil, err
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return NewProjectsRepository(sqlxDB), db, mock, nil
}

func newMockTasksRepo() (projects.TasksRepository, *sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, nil, err
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return NewTasksRepository(sqlxDB), db, mock, nil
}
