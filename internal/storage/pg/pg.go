package pg

import (
	"Habr-comments-server/internal/config"
	"Habr-comments-server/internal/models"
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
	err := s.db.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}

//TODO: Return list of posts @SysteamPost
// Return post by id and return comments @SysteamPost
// Return comments in comments
// Write block comments by id post @SysteamPost
// Add comment on parent or id comment

// Добавляет пост
func (s *Storage) SavePost(ctx context.Context) (models.Post, error) {
	const op = "storage.post.New"

	panic("storage.post.New: implement me")
}

// Возвращает старший комментарий
func (s *Storage) ParentComments(ctx context.Context, id int) ([]models.Comment, error) {
	const op = "storage.post.ParentComments"
	panic("storage.post.ParentComments: implement me")
}

// Возвращает младший комментарий
func (s *Storage) CommentsWithParent(ctx context.Context, id int) (models.Comment, error) {
	const op = "storage.post.CommentsWithParent"
	panic("storage.post.CommentsWithParent: implement me")
}

// Блокировка комментариев на посте
func (s *Storage) ModeBlockComment(ctx context.Context, id int) (models.Comment, error) {
	const op = "storage.post.ModeBlockComment"
	panic("storage.post.ModeBlockComment: implement me")
}

// Добавляет комментарий
func (s *Storage) SaveComment(ctx context.Context) (models.Comment, error) {
	const op = "storage.post.AddComment"
	panic("storage.post.AddComment: implement me")
}

// Добавляет вложенный комментарий
func (s *Storage) SaveJuniorComment(ctx context.Context, id int) error {
	const op = "storage.post.AddJuniorComment"
	panic("storage.post.AddComment: implement me")
}
