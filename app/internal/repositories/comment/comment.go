package comment

import (
	"app/internal/commons"
	"app/internal/entities"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name=CommentRepository --output=mocks --outpkg=mocks
type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entities.Comment) error
	GetCommentsByPostId(ctx context.Context, postId uint, limit, offset int) ([]entities.Comment, error)
}

type commentRepo struct {
	db             *gorm.DB
	ContextTimeout time.Duration
}

func NewCommentRepository(db *gorm.DB, timeout time.Duration) CommentRepository {
	return &commentRepo{
		db:             db,
		ContextTimeout: timeout,
	}
}

// CreateComment inserts a new comment into the database
func (r *commentRepo) CreateComment(ctx context.Context, comment *entities.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	err := r.db.WithContext(ctx).Create(comment).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return commons.ErrTimeout
		}
		return err
	}
	return nil
}

// GetCommentsByPostId returns comments for a given post ID with pagination
func (r *commentRepo) GetCommentsByPostId(ctx context.Context, postId uint, limit, offset int) ([]entities.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ContextTimeout)
	defer cancel()

	var comments []entities.Comment
	// please order the comments by created_at in descending order
	err := r.db.WithContext(ctx).Where("post_id = ?", postId).Limit(limit).Offset(offset).Preload("Author").Order("created_at desc").Find(&comments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commons.ErrNotFound
		}

		if ctx.Err() == context.DeadlineExceeded {
			return nil, commons.ErrTimeout
		}
		return nil, err
	}
	return comments, nil
}
