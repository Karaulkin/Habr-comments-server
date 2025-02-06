package tservice

import (
	"Habr-comments-server/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) GetComments(ctx context.Context, postID int) ([]models.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func (m *MockCommentService) GetChildComments(ctx context.Context, parentID int) ([]models.Comment, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]models.Comment), args.Error(1)
}

func (m *MockCommentService) CreateComment(ctx context.Context, postID, authorID int, parentID *int, content string) (int, error) {
	args := m.Called(ctx, postID, authorID, parentID, content)
	//log.Printf("Переданный ctx: %v", ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockCommentService) GetChildCommentsByParentID(ctx context.Context, parentIDs []int) ([][]*models.Comment, error) {
	args := m.Called(ctx, parentIDs)
	return args.Get(0).([][]*models.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsByPostID(ctx context.Context, postIDs []int) ([][]*models.Comment, error) {
	args := m.Called(ctx, postIDs)
	return args.Get(0).([][]*models.Comment), args.Error(1)
}
