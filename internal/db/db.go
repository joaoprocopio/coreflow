package db

import (
	"context"
	"convey/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB = pgx.Conn

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbDatabase,
	)
	db, err := pgx.Connect(ctx, dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
