package service

import (
	"Habr-comments-server/internal/models"
	"context"
)

type Service struct {
	PostService    PostService
	CommentService CommentService
}

// Конструктор Service
func NewService(postService PostService, commentService CommentService) *Service {
	return &Service{
		PostService:    postService,
		CommentService: commentService,
	}
}

type PostService interface {
	GetPost(ctx context.Context, id int) (models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
	CreatePost(ctx context.Context, authorId int, title, content string, allowComments bool) (int, error)
	BlockComments(ctx context.Context, id int) error
	GetUsersByID(ctx context.Context, ids []int) ([]*models.User, error)
}

type CommentService interface {
	GetComments(ctx context.Context, postID int) ([]models.Comment, error)
	GetChildComments(ctx context.Context, parentID int) ([]models.Comment, error)
	CreateComment(ctx context.Context, postID, authorID int, parentID *int, content string) (int, error)
	GetChildCommentsByParentID(ctx context.Context, parentIDs []int) ([][]*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postIDs []int) ([][]*models.Comment, error)
}
