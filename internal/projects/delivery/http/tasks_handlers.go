package http

import (
	"context"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/internal/projects"
	"github.com/armanokka/time_tracker/pkg/httpErrors"
	"github.com/armanokka/time_tracker/pkg/logger"
	"github.com/armanokka/time_tracker/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tasksHandlers struct {
	tasksUC projects.TasksUseCase
	log     logger.Logger
	tracer  trace.Tracer
}

// Get godoc
// @Summary      Get project tasks
// @Description  Get project tasks
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  []models.Task
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks [get]
func (h tasksHandlers) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Get")
		defer span.End()

		projectID := c.GetInt64("project_id")

		results, err := h.tasksUC.Get(ctx, projectID)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, results)
	}
}

// Create godoc
// @Summary      Create project task
// @Description  Create project task
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param		 taskBody body  models.Task true "task to be created"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks [post]
func (h tasksHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Create")
		defer span.End()

		task := &models.Task{}
		if err := utils.ReadRequest(c, task); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		task.ProjectID = c.GetInt64("project_id")

		task, err := h.tasksUC.Create(ctx, task)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, task)
	}
}

// Update godoc
// @Summary      Update project task
// @Description  Update project task
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param		 updatedTaskBody body  models.Task true "updates to the task"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id} [patch]
func (h tasksHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Update")
		defer span.End()

		task := &models.Task{}
		if err := utils.ReadRequest(c, task); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		task.ID = c.GetInt64("task_id")

		task, err := h.tasksUC.Update(ctx, task)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		c.JSON(200, task)
	}
}

// Delete godoc
// @Summary      Delete project task
// @Description  Delete project task
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id} [delete]
func (h tasksHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Update")
		defer span.End()

		taskID := c.GetInt64("task_id")

		if err := h.tasksUC.Delete(ctx, taskID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

// Start godoc
// @Summary      Start doing project task
// @Description  Start doing project task
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id}/start [post]
func (h tasksHandlers) Start() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Start")
		defer span.End()

		user := c.MustGet("user").(*models.User)
		taskID := c.GetInt64("task_id")

		if err := h.tasksUC.Start(ctx, taskID, user.ID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		c.JSON(200, utils.Response{Ok: true})
	}
}

// Stop godoc
// @Summary      Stop doing project task
// @Description  Stop doing project task
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id}/stop [post]
func (h tasksHandlers) Stop() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.Stop")
		defer span.End()

		user := c.MustGet("user").(*models.User)
		taskID := c.GetInt64("task_id")

		if err := h.tasksUC.Stop(ctx, taskID, user.ID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		c.JSON(200, utils.Response{Ok: true})
	}
}

// GetMembers godoc
// @Summary      Get task executors
// @Description  Get task executors
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  []models.User
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id}/users [get]
func (h tasksHandlers) GetMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.GetMembers")
		defer span.End()

		taskID := c.GetInt64("task_id")

		members, err := h.tasksUC.GetMembers(ctx, taskID)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, members)
	}
}

// AddMember godoc
// @Summary      Add task executor
// @Description  Add task executor
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Param		 addExecutorRequestBody body  http.AddTaskMemberRequest true "add executor request body"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id}/users [post]
func (h tasksHandlers) AddMember() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.AddMember")
		defer span.End()

		taskID := c.GetInt64("task_id")

		req := &AddTaskMemberRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		if err := h.tasksUC.AddMember(ctx, taskID, req.UserID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

// DeleteMember godoc
// @Summary      Remove task executor
// @Description  Remove task executor
// @Tags		 tasks
// @Produce      json
// @Param        project_id path string true "project id"
// @Param        task_id path string true "task id"
// @Param        user_id path string true "executor id to be removed"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /projects/{project_id}/tasks/{task_id}/users/{user_id} [delete]
func (h tasksHandlers) DeleteMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "tasksHandlers.RemoveMember")
		defer span.End()

		taskID := c.GetInt64("task_id")
		userID := c.GetInt64("user_id")

		if err := h.tasksUC.DeleteMember(ctx, taskID, userID); err != nil {
			utils.LogResponseError(c, h.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

func NewTasksHandlers(tasksUC projects.TasksUseCase, log logger.Logger) projects.TaskHandlers {
	return tasksHandlers{tasksUC: tasksUC, tracer: otel.GetTracerProvider().Tracer("api"),
		log: log}
}
