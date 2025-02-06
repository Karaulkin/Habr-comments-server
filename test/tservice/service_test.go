package tservice

import (
	"Habr-comments-server/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	s "Habr-comments-server/internal/service"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestGetPost(t *testing.T) {
	mockPostService := new(MockPostService)
	mockCommentService := new(MockCommentService)
	service := s.NewService(mockPostService, mockCommentService)

	ctx := context.Background()
	expectedPost := models.Post{ID: 1, Title: "Test Post", Content: "Test Content", AllowComments: true}

	mockPostService.On("GetPost", ctx, 1).Return(expectedPost, nil)

	post, err := service.PostService.GetPost(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)

	mockPostService.AssertExpectations(t)
}

func TestCreatePost(t *testing.T) {
	mockPostService := new(MockPostService)
	mockCommentService := new(MockCommentService)
	service := s.NewService(mockPostService, mockCommentService)

	ctx := context.Background()
	mockPostService.On("CreatePost", ctx, 1, "Title", "Content", true).Return(1, nil)

	postID, err := service.PostService.CreatePost(ctx, 1, "Title", "Content", true)
	assert.NoError(t, err)
	assert.Equal(t, 1, postID)

	mockPostService.AssertExpectations(t)
}

func TestCreateComment(t *testing.T) {
	mockPostService := new(MockPostService)
	mockCommentService := new(MockCommentService)
	service := s.NewService(mockPostService, mockCommentService)

	ctx := context.Background() // ОДИН раз создаем контекст

	mockCommentService.On("CreateComment", mock.Anything, 1, 2, (*int)(nil), "New Comment").Return(10, nil).Once()

	// Вызываем тестируемую функцию
	commentID, err := service.CommentService.CreateComment(ctx, 1, 2, nil, "New Comment")

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, 10, commentID)

	// Проверяем вызовы моков
	mockPostService.AssertExpectations(t)
	mockCommentService.AssertExpectations(t)
}

func TestBlockComments(t *testing.T) {
	mockPostService := new(MockPostService)
	mockCommentService := new(MockCommentService)
	service := s.NewService(mockPostService, mockCommentService)

	ctx := context.Background()
	mockPostService.On("BlockComments", ctx, 1).Return(nil)

	err := service.PostService.BlockComments(ctx, 1)
	assert.NoError(t, err)

	mockPostService.AssertExpectations(t)
}
