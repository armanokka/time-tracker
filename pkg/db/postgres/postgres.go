package postgres

import (
	"context"
	"fmt"
	"github.com/armanokka/test_task_Effective_mobile/config"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(ctx context.Context, c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DB,
		c.Postgres.Password,
	)

	db, err := sqlx.ConnectContext(ctx, c.Postgres.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
