package postgres

import (
	"context"
	"fmt"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(ctx context.Context, cfg *config.PostgresConfig) (*sqlx.DB, error) {
	// Connecting to database
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DB,
		cfg.Password,
	)

	db, err := sqlx.ConnectContext(ctx, cfg.Driver, dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "postgres.NewPsqlDB.ConnectContext")
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "postgres.NewPsqlDB.Ping")
	}

	// Running migrations
	m, err := migrate.New("file://migrations", fmt.Sprintf("pgx5://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB))
	if err != nil {
		return nil, errors.Wrap(err, "postgres.NewPsqlDB.New")
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, errors.Wrap(err, "postgres.NewPsqlDB.Up")
	}

	return db, nil
}
