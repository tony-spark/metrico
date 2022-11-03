package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DatabaseRepository struct {
	DB *sql.DB
}

func NewPsqlRepository(dsn string) (*DatabaseRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepository{
		DB: db,
	}, nil
}

func (d DatabaseRepository) Check(ctx context.Context) (bool, error) {
	ct, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.DB.PingContext(ct); err != nil {
		return false, err
	}
	return true, nil
}

func (d DatabaseRepository) Close() error {
	return d.DB.Close()
}
