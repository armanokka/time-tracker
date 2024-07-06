package repository

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
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

func TestProjectsRedisRepo_GetProject(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewProjectsRedisRepo(rdb)
	project := getTestProject()

	gotProject, err := repo.GetProject(ctx, project.ID)
	assert.Nil(t, gotProject)
	assert.NotNil(t, err)

	assert.Nil(t, repo.SetProject(ctx, project, 10))
	gotProject, err = repo.GetProject(ctx, project.ID)
	assert.Nil(t, err)
	assert.Equal(t, project, gotProject)
}

func TestProjectsRedisRepo_DeleteProject(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewProjectsRedisRepo(rdb)
	project := getTestProject()

	assert.Nil(t, repo.SetProject(ctx, project, 10))
	gotProject, err := repo.GetProject(ctx, project.ID)
	assert.Nil(t, err)
	assert.Equal(t, project, gotProject)

	assert.Nil(t, repo.DeleteProject(ctx, project.ID))
	gotProject, err = repo.GetProject(ctx, project.ID)
	assert.NotNil(t, err)
	assert.Nil(t, gotProject)
}

func TestProjectsRedisRepo_SetProject(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewProjectsRedisRepo(rdb)
	project := getTestProject()

	assert.Nil(t, repo.SetProject(ctx, project, 10))
	gotProject, err := repo.GetProject(ctx, project.ID)
	assert.Equal(t, project, gotProject)
	assert.Nil(t, err)
}
