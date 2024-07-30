package repository

import (
	"context"
	"database/sql"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type projectsRepo struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewProjectsRepository(db *sqlx.DB) projects.Repository {
	return projectsRepo{db: db, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (c projectsRepo) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.CreateProject")
	defer span.End()

	var createdProject models.Project
	if err := c.db.QueryRowxContext(ctx, createProjectQuery, project.Name, project.Description,
		project.CreatorID).StructScan(&createdProject); err != nil {
		return nil, err
	}
	if err := c.AddMember(ctx, createdProject.ID, createdProject.CreatorID); err != nil {
		return nil, err
	}
	return &createdProject, nil
}

func (c projectsRepo) GetByID(ctx context.Context, projectID int64) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.GetProjectByID")
	defer span.End()

	project := &models.Project{}
	return project, c.db.QueryRowxContext(ctx, getProjectByIDQuery, projectID).StructScan(project)
}

func (c projectsRepo) Delete(ctx context.Context, projectID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.DeleteProject")
	defer span.End()

	result, err := c.db.ExecContext(ctx, deleteProjectQuery, projectID)
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

func (c projectsRepo) Update(ctx context.Context, updatedProject *models.Project) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.UpdateProject")
	defer span.End()

	return updatedProject, c.db.QueryRowxContext(ctx, updateProjectQuery, updatedProject.Name, updatedProject.Description,
		updatedProject.CreatorID, updatedProject.ID).StructScan(updatedProject)
}

func (c projectsRepo) IsMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.IsProjectMember")
	defer span.End()

	result, err := c.db.ExecContext(ctx, isProjectMemberQuery, projectID, userID)
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

func (c projectsRepo) IsOwner(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.IsProjectMember")
	defer span.End()

	result, err := c.db.ExecContext(ctx, isProjectOwnerQuery, projectID, userID)
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

func (c projectsRepo) AddMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.AddProjectMember")
	defer span.End()

	_, err := c.db.ExecContext(ctx, addProjectMemberQuery, projectID, userID)
	return err
}

func (c projectsRepo) RemoveMember(ctx context.Context, projectID, userID int64) error {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.RemoveProjectMember")
	defer span.End()

	result, err := c.db.ExecContext(ctx, removeProjectMemberQuery, projectID, userID)
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

func (c projectsRepo) GetMembers(ctx context.Context, projectID int64) ([]*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.GetMembers")
	defer span.End()

	var membersCount int
	if err := c.db.GetContext(ctx, &membersCount, getProjectMembersCount, projectID); err != nil {
		return nil, err
	}

	rows, err := c.db.QueryxContext(ctx, getProjectMembers, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0, membersCount)

	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (c projectsRepo) GetMemberProductivity(ctx context.Context, projectID, userID int64) ([]models.UserProductivity, error) {
	ctx, span := c.tracer.Start(ctx, "projectsRepo.GetMemberProductivity")
	defer span.End()

	rows, err := c.db.QueryxContext(ctx, getProjectMemberProductivityQuery, projectID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productivity = make([]models.UserProductivity, 0, 10)
	for rows.Next() {
		result := struct {
			TaskID       int64   `db:"task_id"`
			TotalSeconds float64 `db:"total_seconds"`
		}{}
		if err = rows.StructScan(&result); err != nil {
			return nil, err
		}
		hours, minutes := getHoursMinutes(int(result.TotalSeconds))

		productivity = append(productivity, models.UserProductivity{
			TaskID:       result.TaskID,
			SpentHours:   hours,
			SpentMinutes: minutes,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productivity, nil
}

func getHoursMinutes(totalSeconds int) (hours, minutes int) {
	hours = totalSeconds / 3600
	minutes = (totalSeconds % 3600) / 60
	return
}
