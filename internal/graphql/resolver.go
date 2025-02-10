package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"Habr-comments-server/internal/graphql/loaders"
	"Habr-comments-server/internal/models"
	"Habr-comments-server/internal/service"
	"context"
	"fmt"
	"strconv"
	"time"
)

type Resolver struct {
	Service *service.Service
	Loaders *loaders.Loaders
}

// Post is the resolver for the post field.
func (r *commentResolver) Post(ctx context.Context, obj *models.Comment) (*models.Post, error) {
	post, err := r.Service.PostService.GetPost(ctx, obj.PostId)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	return r.Loaders.UserLoader.Load(obj.AuthorId)
}

// Parent is the resolver for the parent field.
func (r *commentResolver) Parent(ctx context.Context, obj *models.Comment) (*models.Comment, error) {
	if obj.ParentId == nil {
		return nil, nil
	}

	comments, err := r.Loaders.CommentLoader.Load(*obj.ParentId)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, nil // Родительского комментария нет
	}

	return comments[0], nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// Children is the resolver for the children field.
func (r *commentResolver) Children(ctx context.Context, obj *models.Comment, limit *int, offset *int) ([]*models.Comment, error) {
	comments, err := r.Service.CommentService.GetChildComments(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	// Реализуем пагинацию
	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil && start+*limit < end {
		end = start + *limit
	}

	if start > end {
		return []*models.Comment{}, nil
	}

	commentPtrs := make([]*models.Comment, end-start)
	for i := range commentPtrs {
		commentPtrs[i] = &comments[start+i]
	}

	return commentPtrs, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, authorID string, title string, content string, allowComments bool) (*models.Post, error) {
	var err error
	var authorIdInt, postID int
	var post models.Post

	authorIdInt, err = strconv.Atoi(authorID)
	if err != nil {
		return nil, err
	}

	postID, err = r.Service.PostService.CreatePost(ctx, authorIdInt, title, content, allowComments)
	if err != nil {
		return nil, err
	}

	post, err = r.Service.PostService.GetPost(ctx, postID)
	return &post, err
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, authorID string, parentID *string, content string) (*models.Comment, error) {
	var err error
	var authorIdInt, commentID, postIdInt int

	// Конвертация строковых параметров в int
	postIdInt, err = strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	authorIdInt, err = strconv.Atoi(authorID)
	if err != nil {
		return nil, fmt.Errorf("invalid author ID: %w", err)
	}

	var parentIdInt *int
	if parentID != nil {
		pID, err := strconv.Atoi(*parentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID: %w", err)
		}
		parentIdInt = &pID
	}

	post, err := r.Service.PostService.GetPost(ctx, postIdInt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	if !post.AllowComments {
		return nil, fmt.Errorf("comments are disabled for this post")
	}

	commentID, err = r.Service.CommentService.CreateComment(ctx, postIdInt, authorIdInt, parentIdInt, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	comments, err := r.Service.CommentService.GetComments(ctx, postIdInt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	for _, comment := range comments {
		if comment.ID == commentID {
			return &comment, nil
		}
	}

	return nil, fmt.Errorf("comment not found after creation")
}

// BlockComments is the resolver for the blockComments field.
func (r *mutationResolver) BlockComments(ctx context.Context, postID string) (*models.Post, error) {
	var err error
	var postIdInt int

	postIdInt, err = strconv.Atoi(postID)
	if err != nil {
		return nil, err
	}

	err = r.Service.PostService.BlockComments(ctx, postIdInt)
	if err != nil {
		return nil, err
	}

	var post models.Post

	post, err = r.Service.PostService.GetPost(ctx, postIdInt)
	return &post, err
}

// Author is the resolver for the author field.
func (r *postResolver) Author(ctx context.Context, obj *models.Post) (*models.User, error) {
	return r.Loaders.UserLoader.Load(obj.AuthorId)
}

// CreatedAt is the resolver for the createdAt field.
func (r *postResolver) CreatedAt(ctx context.Context, obj *models.Post) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *models.Post, limit *int, offset *int) ([]*models.Comment, error) {
	comments, err := r.Loaders.CommentLoader.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	// Реализуем пагинацию
	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil && start+*limit < end {
		end = start + *limit
	}

	if start > end { // Проверяем, не ушли ли мы за границы массива
		return []*models.Comment{}, nil
	}

	return comments[start:end], nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	posts, err := r.Service.PostService.GetPosts(ctx, *limit, *offset)
	if err != nil {
		return nil, err
	}

	postPtrs := make([]*models.Post, len(posts))
	for i := range posts {
		postPtrs[i] = &posts[i]
	}
	return postPtrs, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	postIdInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	post, err := r.Service.PostService.GetPost(ctx, postIdInt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, parentID string, limit *int, offset *int) ([]*models.Comment, error) {
	parentIdInt, err := strconv.Atoi(parentID)
	if err != nil {
		return nil, err
	}

	comments, err := r.Service.CommentService.GetChildComments(ctx, parentIdInt)
	if err != nil {
		return nil, err
	}

	// Применяем пагинацию (если limit и offset заданы)
	start := 0
	if offset != nil {
		start = *offset
	}
	end := len(comments)
	if limit != nil && start+*limit < end {
		end = start + *limit
	}

	// Преобразуем []models.Comment -> []*models.Comment
	commentPtrs := make([]*models.Comment, end-start)
	for i := range commentPtrs {
		commentPtrs[i] = &comments[start+i]
	}

	return commentPtrs, nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
