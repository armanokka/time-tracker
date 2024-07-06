package auth

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int64) error
	SearchUsers(ctx context.Context, req *utils.UsersQuery) (utils.UsersQueryResponse, error)
}
