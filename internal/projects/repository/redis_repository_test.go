package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

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
