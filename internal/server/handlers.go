package server

import (
	authHttp "github.com/armanokka/test_task_Effective_mobile/internal/auth/delivery/http"
	authRepo "github.com/armanokka/test_task_Effective_mobile/internal/auth/repository"
	authUc "github.com/armanokka/test_task_Effective_mobile/internal/auth/usecase"
	"github.com/armanokka/test_task_Effective_mobile/internal/middleware"
	projectsHttp "github.com/armanokka/test_task_Effective_mobile/internal/projects/delivery/http"
	projectsRepo "github.com/armanokka/test_task_Effective_mobile/internal/projects/repository"
	projectsUc "github.com/armanokka/test_task_Effective_mobile/internal/projects/usecase"
	"github.com/gin-gonic/gin"
)

func (s Server) MapHandlers(c *gin.RouterGroup) {
	aRepo := authRepo.NewAuthRepository(s.db)                                  // auth repository
	aRedisRepo := authRepo.NewAuthRedisRepo(s.rdb)                             // auth redis repository
	aUseCase := authUc.NewAuthUseCase(s.cfg.Server, aRepo, aRedisRepo)         // auth use case
	authHandlers := authHttp.NewAuthHandlers(s.cfg.Server, aUseCase, s.logger) // auth handlers

	projRepo := projectsRepo.NewProjectsRepository(s.db)      // projects repository
	projRedisRepo := projectsRepo.NewProjectsRedisRepo(s.rdb) // projects redis repository
	tasksRepo := projectsRepo.NewTasksRepository(s.db)        // tasks repository

	projectsUC := projectsUc.NewProjectsUseCase(projRepo, projRedisRepo) // projects use case
	tasksUC := projectsUc.NewTasksUseCase(tasksRepo)                     // tasks use case

	projectsHandlers := projectsHttp.NewProjectsHandlers(s.cfg.Server, projectsUC, s.logger) // projects handlers
	tasksHandlers := projectsHttp.NewTasksHandlers(tasksUC, s.logger)                        // tasks handlers

	mw := middleware.NewMiddlewareManager(s.cfg.Server, []string{"*"}, s.logger, aUseCase, projectsUC, tasksUC)
	s.router.Use(mw.ParsePathParametersMiddleware())

	authHttp.MapAuthRoutes(c.Group("/users"), authHandlers, mw)
	projectsHttp.MapProjectsTasksRoutes(c.Group("/projects"), projectsHandlers, tasksHandlers, mw)
}
