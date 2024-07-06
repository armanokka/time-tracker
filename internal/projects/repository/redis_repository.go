package repository

import (
	"context"
	"encoding/json"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/internal/projects"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type projectsRedisRepo struct {
	rdb    *redis.Client
	tracer trace.Tracer
}

func NewProjectsRedisRepo(rdb *redis.Client) projects.RedisRepository {
	return projectsRedisRepo{rdb: rdb, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (c projectsRedisRepo) SetProject(ctx context.Context, project *models.Project, seconds int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRedisRepo.SaveProject")
	defer span.End()

	jsonProject, err := json.Marshal(project)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, "project:"+strconv.FormatInt(project.ID, 10), string(jsonProject), time.Duration(seconds)*time.Second).Err()
}
func (c projectsRedisRepo) GetProject(ctx context.Context, projectID int64) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRedisRepo.GetProjectByID")
	defer span.End()

	cmd := c.rdb.Get(ctx, "project:"+strconv.FormatInt(projectID, 10))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var user models.Project
	if err := json.Unmarshal([]byte(cmd.Val()), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (c projectsRedisRepo) DeleteProject(ctx context.Context, projectID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRedisRepo.DeleteProject")
	defer span.End()

	return c.rdb.Del(ctx, "project:"+strconv.FormatInt(projectID, 10)).Err()
}
