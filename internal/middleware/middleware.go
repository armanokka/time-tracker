package middleware

import (
	"database/sql"
	"errors"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/armanokka/test_task_Effective_mobile/internal/auth"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/internal/projects"
	"github.com/armanokka/test_task_Effective_mobile/pkg/httpErrors"
	"github.com/armanokka/test_task_Effective_mobile/pkg/logger"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type Manager struct {
	cfg        config.ServerConfig
	origins    []string
	log        logger.Logger
	tracer     trace.Tracer
	authUC     auth.UseCase
	projectsUC projects.UseCase
	tasksUC    projects.TasksUseCase
}

func NewMiddlewareManager(cfg config.ServerConfig, origins []string, log logger.Logger,
	authUC auth.UseCase, projectsUC projects.UseCase, tasksUC projects.TasksUseCase) Manager {
	return Manager{
		cfg:        cfg,
		origins:    origins,
		log:        log,
		authUC:     authUC,
		projectsUC: projectsUC,
		tracer:     otel.GetTracerProvider().Tracer("api"),
		tasksUC:    tasksUC,
	}
}

func (m Manager) AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Access-Token")
		if token == "" {
			m.log.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				requestid.Get(c),
				"empty X-Access-Token",
			)
			c.AbortWithStatusJSON(403, httpErrors.NewUnauthorizedError("empty X-Access-Token"))
			return
		}

		if err := m.validateJWTToken(c, token); err != nil {
			c.AbortWithStatusJSON(403, httpErrors.NewUnauthorizedError("wrong X-Access-Token"))
		}
	}
}

func (m Manager) ParsePathParametersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("project_id") != "" {
			projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Param(project_id) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid project_id")
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
			c.Set("project_id", projectID)
		}
		if c.Param("user_id") != "" {
			userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Param(user_id) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid user_id")
				c.AbortWithStatusJSON(http.StatusBadRequest, httpErrors.NewBadRequestError(httpErrors.BadRequest))
				return
			}
			c.Set("user_id", userID)
		}
		if c.Param("task_id") != "" {
			taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Param(task_id) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid task_id")
				c.AbortWithStatusJSON(http.StatusBadRequest, httpErrors.NewBadRequestError(httpErrors.BadRequest))
				return
			}
			c.Set("task_id", taskID)
		}
	}
}

func (m Manager) MemberOrOwnerOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		if user.Admin {
			return // admins can do everything
		}

		if projectID := c.GetInt64("project_id"); projectID != 0 {
			// Checking if user is the owner of the project
			err := m.projectsUC.IsOwner(c, projectID, user.ID)
			if err == nil {
				return // Since user is owner of the project, we don't need to check his membership in project/tasks
			}
			if !errors.Is(err, sql.ErrNoRows) {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			}

			// Checking if user is a member of the project
			if err := m.projectsUC.IsMember(c, projectID, user.ID); err != nil {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
		}

		if taskID := c.GetInt64("task_id"); taskID != 0 {
			// Checking if user is a member of the task
			if err := m.tasksUC.IsMember(c, taskID, user.ID); err != nil {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
		}

		if userID := c.GetInt64("user_id"); userID != 0 {
			// Verifying that user can edit himself only
			if user.ID != userID {
				err := httpErrors.NewForbiddenError("not enough permissions")
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
		}
	}
}

func (m Manager) OwnerOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		if user.Admin {
			return // admins can do everything
		}

		projectID := c.GetInt64("project_id")
		if projectID == 0 {
			return
		}
		// Checking whether user is the owner of the project
		if err := m.projectsUC.IsOwner(c, projectID, user.ID); err != nil {
			utils.LogResponseError(c, m.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		return
	}
}
