package pg

import (
	"Habr-comments-server/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(dbCfg config.DB) (*Storage, error) {
	const op = "storage.pg.New"

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database) // postgresql://user:password@host:port/database

	pool, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create pool: %w", op, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) Stop(ctx context.Context) error {
	const op = "storage.pg.Stop"

	err := s.db.Close(ctx)

	if err != nil {
		return err
	}

	return nil
}

//TODO: Return list of posts @SysteamPost
// Return post by id and return comments @SysteamPost
// Write block comments by id post @SysteamPost
// Add comment on parent or id comment
