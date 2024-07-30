package auth

import (
	"context"
	"github.com/armanokka/time_tracker/internal/models"
)

type RedisRepository interface {
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	SetUser(ctx context.Context, user *models.User, seconds int) error
	DeleteUser(ctx context.Context, userID int64) error
}
