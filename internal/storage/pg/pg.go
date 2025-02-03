package pg

import (
	"Habr-comments-server/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(dbCfg config.DB) (*Storage, error) {
	const op = "storage.pg.New"

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database) // postgresql://user:password@host:port/database

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create pool: %w", op, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) Stop() error {
	const op = "storage.pg.Stop"
	if s.db != nil {
		s.db.Close()
		return nil
	}
	return fmt.Errorf("%s: err db closed", op)
}
