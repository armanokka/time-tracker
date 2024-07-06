package http

import (
	"github.com/armanokka/test_task_Effective_mobile/internal/auth"
	"github.com/armanokka/test_task_Effective_mobile/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers, mw middleware.Manager) {
	authGroup.Use(mw.ParsePathParametersMiddleware())
	authGroup.GET("/", mw.AuthJWTMiddleware(), h.SearchUsers())
	authGroup.POST("/", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.GET("/:user_id", mw.AuthJWTMiddleware(), h.GetUserByID())
	authGroup.PATCH("/:user_id", mw.AuthJWTMiddleware(), h.Update())
	authGroup.DELETE("/:user_id", mw.AuthJWTMiddleware(), mw.OwnerOrAdminMiddleware(), h.Delete())
}
