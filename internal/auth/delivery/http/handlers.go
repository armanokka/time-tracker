package http

import (
	"context"
	"github.com/armanokka/time_tracker/config"
	"github.com/armanokka/time_tracker/internal/auth"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/pkg/httpErrors"
	"github.com/armanokka/time_tracker/pkg/logger"
	"github.com/armanokka/time_tracker/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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

		paramUserID := c.GetInt64("user_id")

		user, err := a.authUC.GetByID(ctx, userID)
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
// @Param        user_id path string true "user id"
// @Param		 updateBody body  http.UpdateUserRequest true "new user info"s
// @Param        X-Access-Token header string true "Token that you get after authorization/registration"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httpErrors.RestError
// @Failure      500  {object}  httpErrors.RestError
// @Router       /users/{user_id} [patch]
func (a authHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.Start(c.MustGet(utils.UserCtxKey).(context.Context), "authHandlers.Update")
		defer span.End()

		userID := c.GetInt64("user_id")

		update := &UpdateUserRequest{}
		if err := utils.ReadRequest(c, update); err != nil {
			utils.LogResponseError(c, a.log, err)
			c.AbortWithStatusJSON(httpErrors.ErrorResponse(err))
			return
		}
		update.ID = userID

		updatedUser, err := a.authUC.Update(ctx, update)
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

		userID := c.GetInt64("user_id")

		if err := a.authUC.Delete(ctx, userID); err != nil {
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
// @Param		 min_id path integer false "min id"
// @Param		 max_id path integer false "max id"
// @Param		 email path string false "email"
// @Param		 name path string false "name"
// @Param		 surname path string false "surname"
// @Param		 patronymic path string false "patronymic"
// @Param		 address path string false "address"
// @Param		 limit path integer false "results amount limit"
// @Param		 page path integer  false "page"
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
