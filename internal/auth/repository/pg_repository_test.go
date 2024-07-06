package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getTestUser() *models.User {
	patronymic := "Ivanovich"
	return &models.User{
		Email:          "test@test.test",
		Password:       "password",
		Name:           "Ivan",
		Surname:        "Ivanov",
		Patronymic:     &patronymic,
		Address:        "USA, LA",
		Admin:          true,
		PassportNumber: 1234,
		PassportSeries: 123456,
	}
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
		WithArgs(&user.Email, sqlmock.AnyArg(), &user.Name, &user.Surname, &user.Patronymic, &user.Address,
			&user.PassportNumber, &user.PassportSeries).
		WillReturnRows(
			sqlmock.NewRows(user.Columns()).
				AddRow(user.Email, user.Password, user.Name, user.Surname, user.Patronymic,
					user.Address, user.Admin, user.PassportNumber, user.PassportSeries),
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

	mock.ExpectQuery(selectUserByIDQuery).WithArgs(user.ID).WillReturnRows(
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
