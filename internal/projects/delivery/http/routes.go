package http

import (
	"github.com/armanokka/test_task_Effective_mobile/internal/middleware"
	"github.com/armanokka/test_task_Effective_mobile/internal/projects"
	"github.com/gin-gonic/gin"
)

func MapRoutes(projectsGroup *gin.RouterGroup, project projects.Handlers, task projects.TaskHandlers, mw middleware.Manager) {
	projectsGroup.Use(mw.AuthJWTMiddleware())
	projectsGroup.POST("/", project.Create())
	projectsGroup.GET("/:project_id", mw.OwnerOrAdminMiddleware(), project.GetByID())
	projectsGroup.PATCH("/:project_id", mw.OwnerOrAdminMiddleware(), project.Update())
	projectsGroup.DELETE("/:project_id", mw.OwnerOrAdminMiddleware(), project.Delete())

	projectsGroup.GET("/:project_id/users", mw.OwnerOrAdminMiddleware(), project.GetMembers())
	projectsGroup.POST("/:project_id/users", mw.OwnerOrAdminMiddleware(), project.AddMember())
	projectsGroup.GET("/:project_id/users/:user_id", mw.OwnerOrAdminMiddleware(), project.GetMemberProductivity())
	projectsGroup.DELETE("/:project_id/users/:user_id", mw.MemberOrOwnerOrAdminMiddleware(), project.RemoveMember())

	tasksGroup := projectsGroup.Group("/:project_id/tasks")
	tasksGroup.Use(mw.MemberOrOwnerOrAdminMiddleware())

	tasksGroup.GET("/", task.Get())
	tasksGroup.POST("/", task.Create())
	//tasksGroup.GET("/:task_id", task.GetByID())
	tasksGroup.PATCH("/:task_id", task.Update())
	tasksGroup.DELETE("/:task_id", task.Delete())

	tasksGroup.POST("/:task_id/start", task.Start())
	tasksGroup.POST("/:task_id/stop", task.Stop())

	tasksGroup.GET("/:task_id/users", task.GetMembers())
	tasksGroup.POST("/:task_id/users", mw.OwnerOrAdminMiddleware(), task.AddMember())
	tasksGroup.DELETE("/:task_id/users/:user_id", mw.OwnerOrAdminMiddleware(), task.DeleteMember())
}
