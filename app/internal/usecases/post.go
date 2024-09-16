package usecases

import (
	"app/internal/commons"
	"app/internal/entities"
	postRepositories "app/internal/repositories/post"
	userRepositories "app/internal/repositories/user"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, post *entities.CreatePostRequest) (*entities.Post, error)
	GetAllPosts(ctx context.Context, limit, offset int) ([]entities.Post, error)
	GetPostByID(ctx context.Context, id uint) (*entities.Post, error)
	UpdatePost(ctx context.Context, post *entities.UpdatePostRequest) (*entities.Post, error)
	DeletePost(ctx context.Context, id uint) error
}

type postUsecase struct {
	postRepo       postRepositories.PostRepository
	userRepo       userRepositories.UserRepository
	contextTimeout time.Duration
}

func NewPostUsecase(post postRepositories.PostRepository, user userRepositories.UserRepository, timeout time.Duration) PostUsecase {
	return &postUsecase{
		postRepo:       post,
		userRepo:       user,
		contextTimeout: timeout,
	}
}

func (u *postUsecase) CreatePost(ctx context.Context, req *entities.CreatePostRequest) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// retrieve user logged in from context
	email := ctx.Value("user").(string)
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	req.AuthorID = user.ID

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	newPost := &entities.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: req.AuthorID,
	}

	err = u.postRepo.CreatePost(ctx, newPost)
	if err != nil {
		return nil, err
	}

	newPost.Author = user

	return newPost, nil
}

func (u *postUsecase) GetAllPosts(ctx context.Context, limit, page int) ([]entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	return u.postRepo.GetAllPosts(ctx, limit, offset)
}

func (u *postUsecase) GetPostByID(ctx context.Context, id uint) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if id == 0 {
		return nil, commons.ErrBadRequest
	}

	return u.postRepo.GetPostById(ctx, id)
}

func (u *postUsecase) UpdatePost(ctx context.Context, req *entities.UpdatePostRequest) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	email := ctx.Value("user").(string)
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	// check if the user is the author of the post
	existingPost, err := u.postRepo.GetPostById(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if existingPost.AuthorID != user.ID {
		return nil, commons.ErrForbidden
	}

	if req.Title != "" {
		existingPost.Title = req.Title
	}
	if req.Content != "" {
		existingPost.Content = req.Content
	}

	err = u.postRepo.UpdatePost(ctx, existingPost)
	if err != nil {
		return nil, err
	}

	return existingPost, nil
}

func (u *postUsecase) DeletePost(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	post, err := u.postRepo.GetPostById(ctx, id)
	if err != nil {
		return err
	}

	email := ctx.Value("user").(string)
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if post.AuthorID != user.ID {
		return commons.ErrForbidden
	}

	return u.postRepo.DeletePost(ctx, id)
}
