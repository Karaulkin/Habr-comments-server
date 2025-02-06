package tservice

import (
	"Habr-comments-server/internal/models"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) GetPost(ctx context.Context, id int) (models.Post, error) {
	args := m.Called(ctx, id) // mock.Anything вместо ctx

	post, ok := args.Get(0).(models.Post)
	if !ok {
		return models.Post{}, fmt.Errorf("unexpected type: %T", args.Get(0))
	}
	return post, args.Error(1)
}

func (m *MockPostService) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *MockPostService) CreatePost(ctx context.Context, authorId int, title, content string, allowComments bool) (int, error) {
	args := m.Called(ctx, authorId, title, content, allowComments)
	return args.Int(0), args.Error(1)
}

func (m *MockPostService) BlockComments(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostService) GetUsersByID(ctx context.Context, ids []int) ([]*models.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*models.User), args.Error(1)
}
