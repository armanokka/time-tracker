package auth

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
)

type UseCase interface {
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	GetByID(ctx context.Context, userID int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int64) error
	SearchUsers(ctx context.Context, req *utils.UsersQuery) (utils.UsersQueryResponse, error)
}
