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
func (s *Storage) SavePost(ctx context.Context, authorId int, title string, content string, allowComment bool) (int, error) {
	const op = "storage.db.SavePost"

	query := `
		INSERT INTO posts (author_id, title, content, allow_comments)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var postID int
	err := s.db.QueryRow(ctx, query, authorId, title, content, allowComment).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert post: %w", op, err)
	}

	return postID, nil
}

// Возвращает пост по его id
func (s *Storage) Post(ctx context.Context, id int) (models.Post, error) {
	const op = "storage.db.Post"

	query := `
	SELECT id, author_id, title, content, allow_comments, created_at
	FROM posts WHERE id = $1;
	`

	var post models.Post
	err := s.db.QueryRow(ctx, query, id).Scan(
		&post.ID, // добавлено
		&post.AuthorId,
		&post.Title,
		&post.Content,
		&post.AllowComment,
		&post.Time,
	)
	if err != nil {
		return models.Post{}, fmt.Errorf("%s: failed to query post: %w", op, err)
	}

	return post, nil
}

// Вернуть все посты
func (s *Storage) Posts(ctx context.Context) ([]models.Post, error) {
	const op = "storage.db.Posts"

	query := `SELECT id, author_id, title, content, allow_comments, created_at FROM posts;`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query posts: %w", op, err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0, 100)
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorId, &post.Title, &post.Content, &post.AllowComment, &post.Time); err != nil {
			return nil, fmt.Errorf("%s: failed to scan post: %w", op, err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return posts, nil
}

// Возвращает комментарий
func (s *Storage) ParentComments(ctx context.Context, postID int) ([]models.Comment, error) {
	const op = "storage.db.ParentComments"

	query := `
	SELECT id, post_id, author_id, parent_id, content, created_at
	FROM comments
	WHERE post_id = $1 AND parent_id IS NULL
	ORDER BY created_at;
	`

	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query parent comments: %w", op, err)
	}
	defer rows.Close()

	comments := make([]models.Comment, 0, 10)
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.Time)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan comment: %w", op, err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Возвращает младший комментарий
func (s *Storage) CommentsWithParent(ctx context.Context, parentID int) ([]models.Comment, error) {
	const op = "storage.db.CommentsWithParent"

	query := `
	SELECT id, post_id, author_id, parent_id, content, created_at
	FROM comments
	WHERE parent_id = $1
	ORDER BY created_at;
	`

	rows, err := s.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query child comments: %w", op, err)
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.Time)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan comment: %w", op, err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Блокировка комментариев на посте
func (s *Storage) ModeBlockComment(ctx context.Context, id int) (int, error) {
	const op = "storage.db.ModeBlockComment"

	query := `UPDATE posts SET allow_comments = FALSE WHERE id = $1 RETURNING id;`

	var postID int
	err := s.db.QueryRow(ctx, query, id).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to block comments: %w", op, err)
	}

	return postID, nil
}

// Добавляет комментарий
func (s *Storage) SaveComment(ctx context.Context, postID int, authorID int, parentID *int, content string) (int, error) {
	const op = "storage.db.SaveComment"

	query := `INSERT INTO comments (post_id, author_id, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id;`

	var commentID int
	var parentIDValue interface{} = nil // если parentID == nil, передаем NULL
	if parentID != nil {
		parentIDValue = *parentID
	}

	err := s.db.QueryRow(ctx, query, postID, authorID, parentIDValue, content).Scan(&commentID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert comment: %w", op, err)
	}

	return commentID, nil
}
