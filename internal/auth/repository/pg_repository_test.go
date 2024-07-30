package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/armanokka/time_tracker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getTestUser() *models.User {
	patronymic := "Ivanovich"
	address := "USA, LA"
	return &models.User{
		Email:      "test@test.test",
		Password:   "password",
		Name:       "Ivan",
		Surname:    "Ivanov",
		Patronymic: &patronymic,
		Address:    &address,
		Admin:      true,
	}
}

func toDriverValue(values []interface{}) ([]driver.Value, error) {
	out := make([]driver.Value, 0, len(values))
	for _, v := range values {
		if !driver.IsValue(v) {
			return nil, fmt.Errorf("toDriverValue: not driver.Value: %v", v)
		}
		out = append(out, v.(driver.Value))
	}
	return out, nil
}

func TestAuthRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	user := getTestUser()

	mock.ExpectQuery(createUserQuery).
		WithArgs(&user.Email, sqlmock.AnyArg(), &user.Name, &user.Surname, &user.Patronymic, &user.Address).WillReturnRows(
		sqlmock.NewRows(user.Columns()).
			AddRow(user.Email, user.Password, user.Name, user.Surname, user.Patronymic,
				user.Address, user.Admin),
	)

	createdUser, err := authRepo.Create(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, user, createdUser)
}

func TestAuthRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	user := getTestUser()

	sql, args, err := squirrel.Select("*").From(pq.QuoteIdentifier("user")).
		Where("id = $1", user.ID).ToSql()
	if err != nil {
		t.Fatal(err)
	}

	values, err := toDriverValue(args)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(sql).WithArgs(values...).WillReturnRows(
		sqlmock.
			NewRows(user.Columns()).
			AddRow(user.Rows()...),
	)

	gotUser, err := authRepo.GetByID(context.Background(), user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, gotUser)
}

func TestAuthRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	user := getTestUser()

	mock.ExpectQuery(selectUserByEmailQuery).WithArgs(user.Email).WillReturnRows(
		sqlmock.
			NewRows(user.Columns()).
			AddRow(user.Rows()...),
	)

	gotUser, err := authRepo.GetByEmail(context.Background(), user.Email)
	assert.Nil(t, err)
	assert.Equal(t, user, gotUser)
}

func TestAuthRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	user := getTestUser()

	mock.ExpectQuery(deleteUserByIDQuery).WithArgs(user.ID).WillReturnRows()

	assert.Nil(t, authRepo.Delete(context.Background(), user.ID))
}
