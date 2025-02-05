package loaders

import (
	"Habr-comments-server/internal/models"
	"Habr-comments-server/internal/service"
	"context"
	"time"
)

type Loaders struct {
	UserLoader         *UserLoader
	CommentLoader      *CommentLoader
	ChildCommentLoader *ChildCommentLoader
}

// Функция для инициализации всех загрузчиков
func NewLoaders(svc *service.Service) *Loaders {
	return &Loaders{
		// Лоадер для пользователей
		UserLoader: &UserLoader{
			wait:     2 * time.Millisecond,
			maxBatch: 100,
			fetch: func(keys []int) ([]*models.User, []error) {
				users, err := svc.PostService.GetUsersByID(context.Background(), keys)
				if err != nil {
					errors := make([]error, len(keys))
					for i := range keys {
						errors[i] = err
					}
					return nil, errors
				}
				return users, nil
			},
		},
		// Лоадер для комментариев по ID постов
		CommentLoader: &CommentLoader{
			wait:     5 * time.Millisecond,
			maxBatch: 50,
			fetch: func(keys []int) ([][]*models.Comment, []error) {
				comments, err := svc.CommentService.GetCommentsByPostID(context.Background(), keys)
				if err != nil {
					errors := make([]error, len(keys))
					for i := range keys {
						errors[i] = err
					}
					return nil, errors
				}
				return comments, nil
			},
		},

		// Лоадер для дочерних комментариев по ID родительских комментариев
		ChildCommentLoader: &ChildCommentLoader{
			wait:     5 * time.Millisecond,
			maxBatch: 50,
			fetch: func(keys []int) ([][]*models.Comment, []error) { // Должно быть [][]*models.Comment
				comments, err := svc.CommentService.GetChildCommentsByParentID(context.Background(), keys)
				if err != nil {
					errors := make([]error, len(keys))
					for i := range keys {
						errors[i] = err
					}
					return nil, errors
				}

				return comments, nil // Убираем flatten, возвращаем [][]*models.Comment
			},
		},
	}
}
