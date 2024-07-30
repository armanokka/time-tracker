package http

import (
	"context"
	"github.com/armanokka/time_tracker/config"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"github.com/armanokka/time_tracker/pkg/httpErrors"
	"github.com/armanokka/time_tracker/pkg/logger"
	"github.com/armanokka/time_tracker/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type projectHandlers struct {
	cfg        config.ServerConfig
	projectsUC projects.UseCase
	log        logger.Logger
	tracer     trace.Tracer
}

func NewProjectsHandlers(cfg config.ServerConfig, projectsUC projects.UseCase, log logger.Logger) projects.Handlers {
	return projectHandlers{
		cfg:        cfg,
		projectsUC: projectsUC,
		log:        log,
		tracer:     otel.GetTracerProvider().Tracer("api"),
	}
}

// GetByID godoc
// @Summary      Get project by ID
// @Description  Get project by ID
// @Tags		 projects
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id} [get]
func (h projectHandlers) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.GetByID")
		defer span.End()

		projectID := c.GetInt64("project_id")

		project, err := h.projectsUC.GetByID(ctx, projectID)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, project)
	}
}

// Create godoc
// @Summary      Create project
// @Description  Create project
// @Tags		 projects
// @Accept       json
// @Produce      json
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Param        project body   models.Project true "Project"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/ [post]
func (h projectHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.Create")
		defer span.End()

		project := &models.Project{}
		if err := utils.ReadRequest(c, project); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		project.CreatorID = c.MustGet("user").(*models.User).ID

		project, err := h.projectsUC.Create(ctx, project)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, project)
	}
}

// Update godoc
// @Summary      Update project info
// @Description  Update project info
// @Tags		 projects
// @Accept       json
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Param        project_updates body models.Project false "Project updates"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id} [patch]
func (h projectHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.Update")
		defer span.End()

		project := &models.Project{}
		if err := utils.ReadRequest(c, project); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		project.ID = c.GetInt64("project_id")

		project, err := h.projectsUC.Update(ctx, project)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, project)
	}
}

// Delete godoc
// @Summary      Delete project
// @Description  Delete project
// @Tags		 projects
// @Accept       json
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id} [delete]
func (h projectHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.Delete")
		defer span.End()

		projectID := c.GetInt64("project_id")

		if err := h.projectsUC.Delete(ctx, projectID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

// AddMember godoc
// @Summary      Add project member
// @Description  Add project member
// @Tags		 projects
// @Accept       json
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Param        user_id body string true "User id of the member you want to invite"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/users [post]
func (h projectHandlers) AddMember() gin.HandlerFunc {
	type Request struct {
		UserID int64 `json:"user_id"`
	}
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.AddMember")
		defer span.End()

		projectID := c.GetInt64("project_id")

		req := &Request{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		if err := h.projectsUC.AddMember(ctx, projectID, req.UserID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		c.JSON(200, utils.Response{Ok: true})
	}
}

// RemoveMember godoc
// @Summary      Remove project member
// @Description  Remove project member
// @Tags		 projects
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        user_id path string true "id of the user you want to remove"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/users/{user_id} [delete]
func (h projectHandlers) RemoveMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.RemoveMember")
		defer span.End()

		projectID := c.GetInt64("project_id")

		userID := c.GetInt64("user_id")

		if err := h.projectsUC.RemoveMember(ctx, projectID, userID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

// GetMembers godoc
// @Summary      Get project members list
// @Description  Get project members list
// @Tags		 projects
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  []models.User
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/users [get]
func (h projectHandlers) GetMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.GetMembers")
		defer span.End()

		projectID := c.GetInt64("project_id")

		members, err := h.projectsUC.GetMembers(ctx, projectID)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, members)
	}
}

// GetMemberProductivity godoc
// @Summary      Get project member productivity
// @Description  Get project member productivity
// @Tags		 projects
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        user_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.UserProductivity
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/users/{user_id} [get]
func (h projectHandlers) GetMemberProductivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "projectHandlers.GetMemberProductivity")
		defer span.End()

		projectID := c.GetInt64("project_id")
		userID := c.GetInt64("user_id")

		productivity, err := h.projectsUC.GetMemberProductivity(ctx, projectID, userID)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		c.JSON(200, productivity)
	}
}
