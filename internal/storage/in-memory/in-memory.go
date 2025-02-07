package in_memory

import (
	"Habr-comments-server/internal/models"
	"Habr-comments-server/internal/service"
	"context"
	"fmt"
	"sync"
	"time"
)

var _ service.PostService = (*InMemoryStorage)(nil)
var _ service.CommentService = (*InMemoryStorage)(nil)

type InMemoryStorage struct {
	posts    map[int]models.Post
	comments map[int][]models.Comment
	users    map[int]models.User
	mu       sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    make(map[int]models.Post),
		comments: make(map[int][]models.Comment),
		users:    make(map[int]models.User),
	}
}

// GetPosts возвращает посты с поддержкой пагинации (limit и offset).
func (s *InMemoryStorage) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Преобразуем карту постов в срез для сортировки
	var posts []models.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	// Пагинация: если offset больше длины среза, возвращаем пустой срез
	if offset >= len(posts) {
		return []models.Post{}, nil
	}

	// Вычисляем конечный индекс с учетом лимита
	end := offset + limit
	if end > len(posts) {
		end = len(posts)
	}

	// Возвращаем срез с постами согласно пагинации
	return posts[offset:end], nil
}

// GetPost возвращает пост по ID.
func (s *InMemoryStorage) GetPost(ctx context.Context, id int) (models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return models.Post{}, fmt.Errorf("post not found")
	}
	return post, nil
}

// CreatePost создает новый пост.
func (s *InMemoryStorage) CreatePost(ctx context.Context, authorId int, title, content string, allowComments bool) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.posts) + 1
	post := models.Post{
		ID:            id,
		AuthorId:      authorId,
		Title:         title,
		Content:       content,
		AllowComments: allowComments,
		CreatedAt:     time.Now(),
	}
	s.posts[id] = post
	return id, nil
}

// BlockComments блокирует комментарии для поста.
func (s *InMemoryStorage) BlockComments(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return fmt.Errorf("post not found")
	}
	post.AllowComments = false
	s.posts[id] = post
	return nil
}

// GetComments возвращает комментарии для поста.
func (s *InMemoryStorage) GetComments(ctx context.Context, postID int) ([]models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments, ok := s.comments[postID]
	if !ok {
		return nil, fmt.Errorf("comments not found for post")
	}
	return comments, nil
}

// GetChildComments возвращает дочерние комментарии для комментария.
func (s *InMemoryStorage) GetChildComments(ctx context.Context, parentID int) ([]models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var childComments []models.Comment
	for _, comment := range s.comments {
		for _, c := range comment {
			if c.ParentId != nil && *c.ParentId == parentID {
				childComments = append(childComments, c)
			}
		}
	}
	return childComments, nil
}

// CreateComment создает новый комментарий.
func (s *InMemoryStorage) CreateComment(ctx context.Context, postID, authorID int, parentID *int, content string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.comments) + 1
	comment := models.Comment{
		ID:        id,
		PostId:    postID,
		AuthorId:  authorID,
		ParentId:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	s.comments[postID] = append(s.comments[postID], comment)
	return id, nil
}

// GetUsersByID возвращает пользователей по ID.
func (s *InMemoryStorage) GetUsersByID(ctx context.Context, ids []int) ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*models.User
	for _, id := range ids {
		user, ok := s.users[id]
		if ok {
			users = append(users, &user)
		}
	}
	return users, nil
}

// GetCommentsByPostID возвращает комментарии для нескольких постов.
func (s *InMemoryStorage) GetCommentsByPostID(ctx context.Context, postIDs []int) ([][]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result [][]*models.Comment
	for _, postID := range postIDs {
		var comments []*models.Comment
		for _, c := range s.comments[postID] {
			comments = append(comments, &c)
		}
		result = append(result, comments)
	}
	return result, nil
}

// GetChildCommentsByParentID возвращает дочерние комментарии для нескольких комментариев.
func (s *InMemoryStorage) GetChildCommentsByParentID(ctx context.Context, parentIDs []int) ([][]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result [][]*models.Comment
	for _, parentID := range parentIDs {
		var comments []*models.Comment
		for _, commentList := range s.comments {
			for _, c := range commentList {
				if c.ParentId != nil && *c.ParentId == parentID {
					comments = append(comments, &c)
				}
			}
		}
		result = append(result, comments)
	}
	return result, nil
}
