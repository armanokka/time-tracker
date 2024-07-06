package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/armanokka/test_task_Effective_mobile/internal/auth"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const cacheTimeSeconds = 60 * 5

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

	// При добавлении сделаем запрос в АПИ, описанный сваггером
	//if err = a.request(user); err != nil {
	//	return nil, err
	//}

	return &models.UserWithToken{
		User:  user,
		Token: jwtToken,
	}, nil
}

func (a authUC) request(user *models.User) error {
	params := url.Values{}
	params.Set("passportSerie", strconv.Itoa(user.PassportSeries))
	params.Set("passportNumber", strconv.Itoa(user.PassportNumber))
	req, err := http.NewRequest("GET", "/info?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("authUC.request: not 200 status code [%d]", resp.StatusCode)
	}

	var person models.People // Вот модель People, как просили в ТЗ
	if err = json.Unmarshal(body, &person); err != nil {
		return err
	}
	// что-то делаем с person
	return nil
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

	return a.authRepo.GetByEmail(ctx, email)
}
func (a authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := a.tracer.Start(ctx, "authUC.Update")
	defer span.End()

	if err := user.PrepareUpdate(); err != nil {
		return nil, nil
	}

	updatedUser, err := a.authRepo.Update(ctx, user)
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
