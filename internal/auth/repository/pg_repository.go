package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/armanokka/time_tracker/internal/auth"
	"github.com/armanokka/time_tracker/internal/auth/delivery/http"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/armanokka/time_tracker/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		user.Surname, user.Patronymic, user.Address).StructScan(user)
}

func (c authRepository) GetByID(ctx context.Context, id int64) (user *models.User, err error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.GetByID")
	defer span.End()

	sql, args, err := squirrel.Select("*").From(pq.QuoteIdentifier("user")).
		Where("id = $1", id).ToSql()
	if err != nil {
		return nil, fmt.Errorf("authRepository.GetByID.Select: %w", err)
	}

	user = &models.User{}
	if err = c.client.QueryRowxContext(ctx, sql, args...).StructScan(user); err != nil {
		return nil, fmt.Errorf("authRepository.GetByID.QueryRowxContext: %w", err)
	}
	return user, nil
}

func (c authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.GetByEmail")
	defer span.End()

	query, args, err := squirrel.Select("*").From(pq.QuoteIdentifier("user")).
		Where("email = $1", email).ToSql()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := c.client.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c authRepository) Delete(ctx context.Context, id int64) error {
	ctx, span := c.tracer.Start(ctx, "authRepository.Delete")
	defer span.End()

	query, args, err := squirrel.Delete(pq.QuoteIdentifier("user")).Where("id = $1", id).ToSql()
	if err != nil {
		return fmt.Errorf("authRepository.Delete.Delete: %w", err)
	}

	_, err = c.client.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("authRepository.Delete.ExecContext: %w", err)
	}
	return nil
}

func (c authRepository) Update(ctx context.Context, updates *http.UpdateUserRequest) (*models.User, error) {
	ctx, span := c.tracer.Start(ctx, "authRepository.Update")
	defer span.End()

	query := squirrel.Update(pq.QuoteIdentifier("user"))

	if updates.Email != nil {
		query = query.Set("email = (?)", *updates.Email)
	}
	if updates.Password != nil {
		query = query.Set("password = (?)", *updates.Password)
	}
	if updates.Name != nil {
		query = query.Set("name = (?)", *updates.Name)
	}
	if updates.Surname != nil {
		query = query.Set("surname = (?)", *updates.Surname)
	}
	if updates.Address != nil {
		query = query.Set("address = (?)", *updates.Address)
	}
	if updates.Patronymic != nil {
		query = query.Set("patronymic = (?)", *updates.Patronymic)
	}

	sql, args, err := query.Where("id = ?", updates.ID).PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, fmt.Errorf("authRepository.Update.ToSql: %w", err)
	}

	var user models.User

	if err = c.client.QueryRowxContext(ctx, sql, args...).StructScan(&user); err != nil {
		return nil, fmt.Errorf("authRepository.Update.QueryRowxContext: %w", err)
	}

	return &user, nil
}

func (c authRepository) SearchUsers(ctx context.Context, query *utils.UsersQuery) (utils.UsersQueryResponse, error) {
	var totalCount int
	if err := c.client.GetContext(ctx, &totalCount, searchUsersCountQuery,
		query.MinID, query.MaxID, query.Email, query.Name, query.Surname, query.Patronymic,
		query.Address); err != nil {
		return utils.UsersQueryResponse{}, err
	}

	rows, err := c.client.QueryxContext(ctx, searchUsersQuery,
		query.MinID, query.MaxID, query.Email, query.Name, query.Surname, query.Patronymic,
		query.Address, query.GetOffset(), query.GetLimit())
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
