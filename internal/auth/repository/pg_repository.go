package repository

import (
	"context"
	"github.com/armanokka/test_task_Effective_mobile/internal/auth"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type authRepository struct {
	client *sqlx.DB
	tracer trace.Tracer
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return authRepository{client: db, tracer: otel.GetTracerProvider().Tracer("api")}
}

func (c authRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.Create")
	defer span.End()

	return user, c.client.QueryRowxContext(ctx, createUserQuery, user.Email, user.Password, user.Name,
		user.Surname, user.Patronymic, user.Address, user.PassportNumber, user.PassportSeries).StructScan(user)
}

func (c authRepository) GetByID(ctx context.Context, id int64) (user *models.User, err error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.GetByID")
	defer span.End()
	user = &models.User{}
	return user, c.client.QueryRowxContext(ctx, selectUserByIDQuery, &id).StructScan(user)
}

func (c authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.GetByEmail")
	defer span.End()

	var user models.User
	if err := c.client.QueryRowxContext(ctx, selectUserByEmailQuery, &email).StructScan(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c authRepository) Delete(ctx context.Context, id int64) error {
	ctx, span := c.tracer.Start(ctx, "authRepository.Delete")
	defer span.End()

	_, err := c.client.ExecContext(ctx, deleteUserByIDQuery, &id)
	return err
}

func (c authRepository) Update(ctx context.Context, updates *models.User) (*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.Update")
	defer span.End()

	return updates, c.client.QueryRowxContext(ctx, updateUserQuery, &updates.Email, &updates.Password, &updates.Name,
		&updates.Surname, &updates.Patronymic, &updates.Address, &updates.PassportNumber,
		&updates.PassportSeries, &updates.ID).StructScan(updates)
}

func (c authRepository) SearchUsers(ctx context.Context, query *utils.UsersQuery) (utils.UsersQueryResponse, error) {
	var totalCount int
	if err := c.client.GetContext(ctx, &totalCount, searchUsersCountQuery,
		query.MinID, query.MaxID, query.Email, query.Name, query.Surname, query.Patronymic,
		query.Address, query.MinPassportNumber, query.MaxPassportNumber,
		query.MinPassportSeries, query.MaxPassportSeries); err != nil {
		return utils.UsersQueryResponse{}, err
	}

	rows, err := c.client.QueryxContext(ctx, searchUsersQuery,
		query.MinID, query.MaxID, query.Email, query.Name, query.Surname, query.Patronymic,
		query.Address, query.MinPassportNumber, query.MaxPassportNumber,
		query.MinPassportSeries, query.MaxPassportSeries, query.GetOffset(), query.GetLimit()) // todo FIX COALESCE query
	if err != nil {
		return utils.UsersQueryResponse{}, err
	}
	defer rows.Close()

	var users = make([]*models.User, 0, query.Limit)
	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			return utils.UsersQueryResponse{}, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return utils.UsersQueryResponse{}, err
	}
	return utils.UsersQueryResponse{
		Users:      users,
		Count:      len(users),
		Page:       query.Page,
		TotalCount: totalCount,
		TotalPages: totalCount / query.Limit,
	}, nil
}
