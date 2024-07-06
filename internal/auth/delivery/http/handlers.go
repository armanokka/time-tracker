package http

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/armanokka/test_task_Effective_mobile/internal/auth"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/httpErrors"
	"github.com/armanokka/test_task_Effective_mobile/pkg/logger"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
)

type authHandlers struct {
	cfg    config.ServerConfig
	authUC auth.UseCase
	log    logger.Logger
	tracer trace.Tracer
}

func NewAuthHandlers(cfg config.ServerConfig, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return authHandlers{cfg: cfg, authUC: authUC, log: log, tracer: otel.GetTracerProvider().Tracer("api")}
}

// Register godoc
// @Summary      Register user
// @Description  Register user
// @Accept       json
// @Produce      json
// @Tags		 auth
// @Param		 user body  models.User true "new user info"
// @Success      200  {object}  models.UserWithToken
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users [post]
func (a authHandlers) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c, "authHandlers.Register")
		defer span.End()

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		createdUser, err := a.authUC.Register(ctx, user)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, createdUser)
	}
}

// Login godoc
// @Summary      Login user
// @Description  Login user
// @Tags		 auth
// @Accept       json
// @Produce      json
// @Param		 searchUserQuery body  http.LoginRequest true "email and password json object"
// @Success      200  {object}  models.UserWithToken
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/login [post]
func (a authHandlers) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c, "authHandlers.Login")
		defer span.End()

		login := &LoginRequest{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		userWithToken, err := a.authUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, userWithToken)
	}
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get user by ID
// @Tags		 auth
// @Produce      json
// @Param        user_id path string true "user id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/{user_id} [get]
func (a authHandlers) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "authHandlers.GetUserByID")
		defer span.End()

		userID := c.MustGet("user").(*models.User).ID

		paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		user, err := a.authUC.GetByID(ctx, paramUserID)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(404, httpErrors.NewNoSuchUserError(err))
			return
		}

		if paramUserID != userID {
			user.Sanitize()
		}

		c.JSON(200, user)
	}
}

// Update godoc
// @Summary      Update user
// @Description  Update user
// @Tags		 auth
// @Accept       json
// @Produce      json
// @Param		 updateBody body  http.UpdateRequest true "new user info"s
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/{user_id} [patch]
func (a authHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "authHandlers.Update")
		defer span.End()

		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		update := &UpdateRequest{}
		if err := utils.ReadRequest(c, update); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		updatedUser, err := a.authUC.Update(ctx, &models.User{
			ID:             userID,
			Email:          update.Email,
			Password:       update.Password,
			Name:           update.Name,
			Surname:        update.Surname,
			Patronymic:     update.Patronymic,
			Address:        update.Address,
			PassportNumber: update.PassportNumber,
			PassportSeries: update.PassportSeries,
		})
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, updatedUser)
	}
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags		 auth
// @Produce      json
// @Param        user_id path string true "user id"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/{user_id} [delete]
func (a authHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "authHandlers.Delete")
		defer span.End()

		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		if err = a.authUC.Delete(ctx, userID); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(500, httpErrors.ParseErrors(err))
			return
		}

		c.JSON(200, utils.Response{Ok: true})
	}
}

// SearchUsers godoc
// @Summary      Search users
// @Description  Search users
// @Tags		 auth
// @Produce      json
// @Param		 searchUserQuery body  utils.UsersQuery true "search query"
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  utils.UsersQueryResponse
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/ [get]
func (a authHandlers) SearchUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "authHandlers.GetAllUsers")
		defer span.End()

		query := &utils.UsersQuery{}
		if err := c.BindQuery(query); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		users, err := a.authUC.SearchUsers(ctx, query)
		if err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}

		c.JSON(200, users)
	}
}
