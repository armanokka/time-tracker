package repository

import (
	"context"
	"encoding/json"
	"github.com/armanokka/time_tracker/internal/auth"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type authRedisRepo struct {
	redisClient *redis.Client
	tracer      trace.Tracer
}

func NewAuthRedisRepo(redisClient *redis.Client) auth.RedisRepository {
	return &authRedisRepo{redisClient: redisClient, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (u authRedisRepo) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	ctx, span := u.tracer.Start(ctx, "authRedisRepo.GetUser")
	defer span.End()

	cmd := u.redisClient.Get(ctx, "user:"+strconv.FormatInt(userID, 10))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	var user models.User
	if err := json.Unmarshal([]byte(cmd.Val()), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u authRedisRepo) SetUser(ctx context.Context, user *models.User, seconds int) error {
	ctx, span := u.tracer.Start(ctx, "authRedisRepo.SetUser")
	defer span.End()

	jsonUser, err := json.Marshal(*user)
	if err != nil {
		return err
	}
	return u.redisClient.Set(ctx, "user:"+strconv.FormatInt(user.ID, 10), string(jsonUser), time.Duration(seconds)*time.Second).Err()
}

func (u authRedisRepo) DeleteUser(ctx context.Context, userID int64) error {
	ctx, span := u.tracer.Start(ctx, "authRedisRepo.DeleteUser")
	defer span.End()

	return u.redisClient.Del(ctx, "user:"+strconv.FormatInt(userID, 10)).Err()
}
