package database

import (
	"context"
	"convey/internal/config"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DB = bun.DB

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	conn := pgdriver.NewConnector(
		pgdriver.WithInsecure(true),
		pgdriver.WithUser(cfg.DbUser),
		pgdriver.WithPassword(cfg.DbPassword),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.DbHost, cfg.DbPort)),
		pgdriver.WithDatabase(cfg.DbDatabase),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
