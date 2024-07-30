package auth

import (
	"context"
	"github.com/armanokka/time_tracker/internal/auth/delivery/http"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/pkg/utils"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, updates *http.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, userID int64) error
	SearchUsers(ctx context.Context, req *utils.UsersQuery) (utils.UsersQueryResponse, error)
}
