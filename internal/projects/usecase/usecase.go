package usecase

import (
	"context"
	"errors"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/internal/projects"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const cacheTimeSeconds = 60 * 5

type projectsUC struct {
	repo      projects.Repository
	redisRepo projects.RedisRepository
	tracer    trace.Tracer
}

func NewProjectsUseCase(repo projects.Repository, redisRepo projects.RedisRepository) projects.UseCase {
	return projectsUC{
		repo:      repo,
		redisRepo: redisRepo,
		tracer:    otel.GetTracerProvider().Tracer("api"),
	}
}

func (c projectsUC) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsUC.CreateProject")
	defer span.End()

	createdProject, err := c.repo.Create(ctx, project)
	if err != nil {
		return nil, err
	}
	if err = c.redisRepo.SetProject(ctx, createdProject, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return createdProject, nil
}

func (c projectsUC) GetByID(ctx context.Context, projectID int64) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsUC.GetProjectByID")
	defer span.End()

	project, err := c.redisRepo.GetProject(ctx, projectID)
	if err == nil {
		return project, nil
	}
	if !errors.Is(err, redis.Nil) {
		return nil, err
	}
	project, err = c.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if err = c.redisRepo.SetProject(ctx, project, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return project, nil
}

func (c projectsUC) Delete(ctx context.Context, projectID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsUC.DeleteProject")
	defer span.End()

	if err := c.redisRepo.DeleteProject(ctx, projectID); err != nil {
		return err
	}
	return c.repo.Delete(ctx, projectID)
}

func (c projectsUC) Update(ctx context.Context, updates *models.Project) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsUC.UpdateProject")
	defer span.End()

	updatedProject, err := c.repo.Update(ctx, updates)
	if err != nil {
		return nil, err
	}
	if err = c.redisRepo.SetProject(ctx, updatedProject, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return updatedProject, nil
}

func (c projectsUC) IsOwner(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsUC.IsProjectOwner")
	defer span.End()

	return c.repo.IsOwner(ctx, projectID, userID)
}

func (c projectsUC) IsMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsUC.IsProjectMember")
	defer span.End()

	return c.repo.IsMember(ctx, projectID, userID)
}

func (c projectsUC) GetMembers(ctx context.Context, projectID int64) ([]*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "projectsUC.GetMembers")
	defer span.End()

	return c.repo.GetMembers(ctx, projectID)
}

func (c projectsUC) AddMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsUC.AddProjectMember")
	defer span.End()

	return c.repo.AddMember(ctx, projectID, userID)
}

func (c projectsUC) RemoveMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsUC.RemoveProjectMember")
	defer span.End()

	return c.repo.RemoveMember(ctx, projectID, userID)
}

func (c projectsUC) GetMemberProductivity(ctx context.Context, projectID, userID int64) ([]models.UserProductivity, error) {
	ctx, span := c.tracer.Start(ctx, "projectsUC.GetMemberProductivity")
	defer span.End()

	return c.repo.GetMemberProductivity(ctx, projectID, userID)
}
