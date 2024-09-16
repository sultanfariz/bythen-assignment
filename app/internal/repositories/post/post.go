package post

import (
	"app/internal/commons"
	"app/internal/entities"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name=PostRepository --output=mocks --outpkg=mocks
type PostRepository interface {
	CreatePost(ctx context.Context, post *entities.Post) error
	GetAllPosts(ctx context.Context, limit, offset int) ([]entities.Post, error)
	GetPostById(ctx context.Context, id uint) (*entities.Post, error)
	UpdatePost(ctx context.Context, post *entities.Post) error
	DeletePost(ctx context.Context, id uint) error
}

type postRepo struct {
	db             *gorm.DB
	ContextTimeout time.Duration
}

func NewPostRepository(db *gorm.DB, timeout time.Duration) PostRepository {
	return &postRepo{
		db:             db,
		ContextTimeout: timeout,
	}
}

// CreatePost inserts a new post into the database
func (r *postRepo) CreatePost(ctx context.Context, post *entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	err := r.db.WithContext(ctx).Create(post).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}
	return nil
}

// GetAllPosts returns all posts with pagination
func (r *postRepo) GetAllPosts(ctx context.Context, limit, offset int) ([]entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	var posts []entities.Post
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Preload("Author").Find(&posts).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, commons.ErrTimeout
		}
		return nil, err
	}
	return posts, nil
}

// GetPostById returns a post by its ID
func (r *postRepo) GetPostById(ctx context.Context, id uint) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	var post entities.Post
	err := r.db.WithContext(ctx).Preload("Author").First(&post, id).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, commons.ErrTimeout
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commons.ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

// UpdatePost updates an existing post
func (r *postRepo) UpdatePost(ctx context.Context, post *entities.Post) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	err := r.db.WithContext(ctx).Save(post).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}
	return nil
}

// DeletePost deletes a post by its ID
func (r *postRepo) DeletePost(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	err := r.db.WithContext(ctx).Delete(&Post{}, id).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}
	return nil
}
