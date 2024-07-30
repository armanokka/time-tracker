package postgres

import (
	"context"
	"fmt"
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

type Config struct {
	User     string
	DB       string
	Password string
	Driver   string
	Host     string
	Port     int
}

func NewPsqlDB(ctx context.Context, cfg *Config) (*sqlx.DB, error) {
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
		return nil, fmt.Errorf("postgres.NewPsqlDB.ConnectContext: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres.NewPsqlDB.Ping: %w", err)
	}

	// Running migrations
	m, err := migrate.New("file://migrations", fmt.Sprintf("pgx5://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB))
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPsqlDB.New: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("postgres.NewPsqlDB.Up: %w", err)
	}

	return db, nil
}
