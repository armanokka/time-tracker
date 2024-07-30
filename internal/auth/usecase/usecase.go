package usecase

import (
	"context"
	"errors"
	"github.com/armanokka/time_tracker/config"
	"github.com/armanokka/time_tracker/internal/auth"
	"github.com/armanokka/time_tracker/internal/auth/delivery/http"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

const cacheTimeSeconds = 60 * 10

type authUC struct {
	cfg       config.ServerConfig
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	tracer    trace.Tracer
}

func NewAuthUseCase(cfg config.ServerConfig, authRepo auth.Repository, redisRepo auth.RedisRepository) auth.UseCase {
	return authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (a authUC) generateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["expires_at"] = time.Now().Add(24 * time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(a.cfg.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a authUC) Login(ctx context.Context, login *models.User) (*models.UserWithToken, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.Login")
	defer span.End()

	user, err := a.authRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}
	if err = user.ComparePassword(login.Password); err != nil {
		return nil, err
	}
	jwtToken, err := a.generateJWT(user)
	if err != nil {
		return nil, err
	}
	if err = a.redisRepo.SetUser(ctx, user, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return &models.UserWithToken{
		User:  user,
		Token: jwtToken,
	}, nil
}

func (a authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.Register")
	defer span.End()

	if err := user.PrepareUpdate(); err != nil {
		return nil, err
	}

	user, err := a.authRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	jwtToken, err := a.generateJWT(user)
	if err != nil {
		return nil, err
	}
	if err = a.redisRepo.SetUser(ctx, user, cacheTimeSeconds); err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User:  user,
		Token: jwtToken,
	}, nil
}

func (a authUC) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.GetByID")
	defer span.End()

	user, err := a.redisRepo.GetUser(ctx, userID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, redis.Nil) {
		return nil, err
	}
	user, err = a.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = a.redisRepo.SetUser(ctx, user, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return user, nil
}

func (a authUC) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.GetByEmail")
	defer span.End()

	return a.authRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
}

func (a authUC) Update(ctx context.Context, updates *http.UpdateRequest) (*models.User, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.Update")
	defer span.End()

	updatedUser, err := a.authRepo.Update(ctx, updates)
	if err != nil {
		return nil, err
	}
	if err = a.redisRepo.SetUser(ctx, updatedUser, cacheTimeSeconds); err != nil {
		return nil, err
	}
	return updatedUser, nil
}
func (a authUC) Delete(ctx context.Context, userID int64) error {
	ctx, span := a.tracer.Start(ctx, "authUC.Delete")
	defer span.End()

	if err := a.redisRepo.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return a.authRepo.Delete(ctx, userID)
}

func (a authUC) SearchUsers(ctx context.Context, req *utils.UsersQuery) (utils.UsersQueryResponse, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.GetAllUsers")
	defer span.End()

	result, err := a.authRepo.SearchUsers(ctx, req)
	if err != nil {
		return utils.UsersQueryResponse{}, err
	}
	for _, user := range result.Users {
		user.Sanitize()
	}
	return result, nil
}
