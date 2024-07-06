package middleware

import (
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

func (m Manager) MemberOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		if c.Param("project_id") != "" {
			projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Get(user) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid project_id")
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
			c.Set("project_id", projectID)
		}

		if c.Param("task_id") != "" {
			taskID, err := strconv.ParseInt(c.Param("task_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Get(user) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid task_id")
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
			c.Set("task_id", taskID)
		}

		if user.Admin {
			return // admins can do everything
		}

		if projectID := c.GetInt64("project_id"); projectID != 0 {
			member, err := m.projectsUC.IsMember(c, projectID, user.ID)
			if err != nil {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, httpErrors.ParseErrors(err))
				return
			}
			if !member {
				m.log.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					requestid.Get(c),
					strconv.FormatInt(user.ID, 10),
					"not a project member",
				)
				c.AbortWithStatusJSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.Forbidden))
				return
			}
		}

		if taskID := c.GetInt64("task_id"); taskID != 0 {
			member, err := m.tasksUC.IsMember(c, taskID, user.ID)
			if err != nil {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, httpErrors.ParseErrors(err))
				return
			}
			if !member {
				m.log.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					requestid.Get(c),
					strconv.FormatInt(user.ID, 10),
					"not a task member",
				)
				c.AbortWithStatusJSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.Forbidden))
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

		// Checking that user is a member of the project
		if c.Param("project_id") != "" {
			projectID, err := strconv.ParseInt(c.Param("project_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Get(user) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid project_id")
				c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
				return
			}
			c.Set("project_id", projectID)

			owner, err := m.projectsUC.IsOwner(c, projectID, user.ID)
			if err != nil {
				utils.LogResponseError(c, m.log, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, httpErrors.ParseErrors(err))
				return
			}
			if !owner {
				member, err := m.projectsUC.IsMember(c, projectID, user.ID)
				if err != nil {
					utils.LogResponseError(c, m.log, err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, httpErrors.ParseErrors(err))
					return
				}
				if !member {
					m.log.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
						requestid.Get(c),
						strconv.FormatInt(user.ID, 10),
						"not a project member",
					)
					c.AbortWithStatusJSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.Forbidden))
					return
				}
			}
		}

		// Checking that user can edit only himself if he's not admin and not a project owner
		if c.Param("user_id") != "" {
			paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
			if err != nil {
				m.log.Errorf("Error c.Get(user) RequestID: %s, ERROR: %s,", requestid.Get(c), "invalid user_id")
				c.AbortWithStatusJSON(http.StatusBadRequest, httpErrors.NewBadRequestError(httpErrors.BadRequest))
				return
			}
			c.Set("user_id", paramUserID)

			if user.ID != paramUserID {
				m.log.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					requestid.Get(c),
					strconv.FormatInt(user.ID, 10),
					"not enough permissions",
				)
				c.AbortWithStatusJSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.Forbidden))
				return
			}
		}
		return
	}
}
