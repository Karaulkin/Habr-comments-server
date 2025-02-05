package pg

import (
	"Habr-comments-server/internal/config"
	"Habr-comments-server/internal/models"
	"Habr-comments-server/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var _ service.PostService = (*Storage)(nil)
var _ service.CommentService = (*Storage)(nil)

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

// Получение всех постов (с поддержкой пагинации)
func (s *Storage) GetPosts(ctx context.Context, limit int, offset int) ([]models.Post, error) {
	const op = "storage.db.GetPosts"

	query := `
	SELECT id, author_id, title, content, allow_comments, created_at
	FROM posts ORDER BY created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query posts: %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	//posts := make([]models.Post, 0, 100)
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.AuthorId,
			&post.Title,
			&post.Content,
			&post.AllowComments,
			&post.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan post: %w", op, err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	//return posts[:len(posts):len(posts)], nil
	return posts, nil
}

// Получение одного поста по ID
func (s *Storage) GetPost(ctx context.Context, idPost int) (models.Post, error) {
	const op = "storage.db.GetPost"

	query := `
	SELECT id, author_id, title, content, allow_comments, created_at
	FROM posts WHERE id = $1;
	`

	var post models.Post
	err := s.db.QueryRow(ctx, query, idPost).Scan(
		&post.ID,
		&post.AuthorId,
		&post.Title,
		&post.Content,
		&post.AllowComments,
		&post.CreatedAt,
	)
	if err != nil {
		return models.Post{}, fmt.Errorf("%s: failed to query post: %w", op, err)
	}

	return post, nil
}

// Создание нового поста
func (s *Storage) CreatePost(ctx context.Context, authorId int, title, content string, allowComments bool) (int, error) {
	const op = "storage.db.CreatePost"

	query := `
		INSERT INTO posts (author_id, title, content, allow_comments)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var postID int
	err := s.db.QueryRow(ctx, query, authorId, title, content, allowComments).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert post: %w", op, err)
	}

	return postID, nil
}

// Блокировка комментариев к посту
func (s *Storage) BlockComments(ctx context.Context, id int) error {
	const op = "storage.db.BlockComments"

	query := `UPDATE posts SET allow_comments = FALSE WHERE id = $1;`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: failed to block comments: %w", op, err)
	}

	return nil
}

// Получение корневых комментариев к посту
func (s *Storage) GetComments(ctx context.Context, postID int) ([]models.Comment, error) {
	const op = "storage.db.GetComments"

	query := `
	SELECT id, post_id, author_id, parent_id, content, created_at
	FROM comments
	WHERE post_id = $1 AND parent_id IS NULL
	ORDER BY created_at;
	`

	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query comments: %w", op, err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan comment: %w", op, err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Получение дочерних комментариев по parentID
func (s *Storage) GetChildComments(ctx context.Context, parentID int) ([]models.Comment, error) {
	const op = "storage.db.GetChildComments"

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

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan comment: %w", op, err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Создание комментария
func (s *Storage) CreateComment(ctx context.Context, postID int, authorID int, parentID *int, content string) (int, error) {
	const op = "storage.db.CreateComment"

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

func (s *Storage) GetUsersByID(ctx context.Context, ids []int) ([]*models.User, error) {
	query := `
		SELECT id, username FROM users WHERE id = ANY($1);
	`

	rows, err := s.db.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("storage.db.GetUsersByID: failed to query users: %w", err)
	}
	defer rows.Close()

	users := make(map[int]*models.User)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, fmt.Errorf("storage.db.GetUsersByID: failed to scan user: %w", err)
		}
		users[user.ID] = &user
	}

	// Собираем пользователей в правильном порядке (как в `ids`)
	result := make([]*models.User, len(ids))
	for i, id := range ids {
		result[i] = users[id] // Если нет в БД, останется `nil`
	}

	return result, nil
}

func (s *Storage) GetCommentsByPostID(ctx context.Context, postIDs []int) ([][]*models.Comment, error) {
	query := `
		SELECT id, post_id, author_id, parent_id, content, created_at
		FROM comments WHERE post_id = ANY($1);
	`

	rows, err := s.db.Query(ctx, query, postIDs)
	if err != nil {
		return nil, fmt.Errorf("storage.db.GetCommentsByPostID: failed to query comments: %w", err)
	}
	defer rows.Close()

	commentMap := make(map[int][]*models.Comment)
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("storage.db.GetCommentsByPostID: failed to scan comment: %w", err)
		}
		commentMap[comment.PostId] = append(commentMap[comment.PostId], &comment)
	}

	result := make([][]*models.Comment, len(postIDs))
	for i, postID := range postIDs {
		result[i] = commentMap[postID]
	}

	return result, nil
}

func (s *Storage) GetChildCommentsByParentID(ctx context.Context, parentIDs []int) ([][]*models.Comment, error) {
	query := `
		SELECT id, post_id, author_id, parent_id, content, created_at
		FROM comments WHERE parent_id = ANY($1);
	`

	rows, err := s.db.Query(ctx, query, parentIDs)
	if err != nil {
		return nil, fmt.Errorf("storage.db.GetChildCommentsByParentID: failed to query child comments: %w", err)
	}
	defer rows.Close()

	commentMap := make(map[int][]*models.Comment)
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.ParentId, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("storage.db.GetChildCommentsByParentID: failed to scan comment: %w", err)
		}
		commentMap[*comment.ParentId] = append(commentMap[*comment.ParentId], &comment)
	}

	result := make([][]*models.Comment, len(parentIDs))
	for i, parentID := range parentIDs {
		result[i] = commentMap[parentID]
	}

	return result, nil
}
