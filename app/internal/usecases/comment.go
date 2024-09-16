package usecases

import (
	"app/internal/commons"
	"app/internal/entities"
	commentRepositories "app/internal/repositories/comment"
	postRepositories "app/internal/repositories/post"
	userRepositories "app/internal/repositories/user"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, comment *entities.CreateCommentRequest) (*entities.Comment, error)
	GetCommentsByPostID(ctx context.Context, postId uint, limit, offset int) ([]entities.Comment, error)
}

type commentUsecase struct {
	commentRepo    commentRepositories.CommentRepository
	postRepo       postRepositories.PostRepository
	userRepo       userRepositories.UserRepository
	contextTimeout time.Duration
}

func NewCommentUsecase(comment commentRepositories.CommentRepository, post postRepositories.PostRepository, user userRepositories.UserRepository, timeout time.Duration) CommentUsecase {
	return &commentUsecase{
		commentRepo:    comment,
		postRepo:       post,
		userRepo:       user,
		contextTimeout: timeout,
	}
}

func (u *commentUsecase) CreateComment(ctx context.Context, req *entities.CreateCommentRequest) (*entities.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// retrieve email from context
	email := ctx.Value("user").(string)
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	req.AuthorID = user.ID

	// Validate request
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	// Check if the post exists
	_, err = u.postRepo.GetPostById(ctx, req.PostID)
	if err != nil {
		return nil, commons.ErrNotFound
	}

	// Create new comment
	newComment := &entities.Comment{
		Content:  req.Content,
		AuthorID: req.AuthorID,
		PostID:   req.PostID,
	}

	err = u.commentRepo.CreateComment(ctx, newComment)
	if err != nil {
		return nil, err
	}

	newComment.Author = user

	return newComment, nil
}

func (u *commentUsecase) GetCommentsByPostID(ctx context.Context, postId uint, limit, page int) ([]entities.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if postId == 0 {
		return nil, commons.ErrBadRequest
	}

	// Check if the post exists
	_, err := u.postRepo.GetPostById(ctx, postId)
	if err != nil {
		return nil, commons.ErrNotFound
	}

	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	return u.commentRepo.GetCommentsByPostId(ctx, postId, limit, offset)
}
