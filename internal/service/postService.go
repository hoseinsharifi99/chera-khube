package service

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostService interface {
	Get(user *model.User, serviceName string) (*model.Post, *model.Adons, int, error)
	GetPostByUser(ctx *gin.Context, serviceName string) (*model.Post, *model.Adons, int, error)
	UpdatePost(post *model.Post) error
}

type postService struct {
	postApiRepo  repository.PostApiRepo
	postDbRepo   repository.PostDbRepo
	addonsDbRepo repository.AdonsDbRepo
	userService  UserService
	config       *helper.ServiceConfig
	logger       *zap.Logger
}

func NewPostService(
	postApiRepo repository.PostApiRepo,
	postDbRepo repository.PostDbRepo,
	addonsDbRepo repository.AdonsDbRepo,
	userService UserService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) PostService {
	return &postService{
		postApiRepo:  postApiRepo,
		postDbRepo:   postDbRepo,
		addonsDbRepo: addonsDbRepo,
		userService:  userService,
		config:       config,
		logger:       logger,
	}
}

func (s postService) Get(user *model.User, serviceName string) (*model.Post, *model.Adons, int, error) {
	if user.PostToken == "" {
		return nil, nil, 0, errors.New("token required")
	}

	post, err := s.postDbRepo.Get(user.PostToken)
	if err != nil {
		return nil, nil, 0, err
	}

	if post == nil {
		post, err = s.postApiRepo.Get(user.PostToken, serviceName)
		if err != nil {
			return nil, nil, 0, err
		}

		post.UserID = user.ID
		post, err = s.postDbRepo.Insert(post)
		if err != nil {
			return nil, nil, 0, err
		}

		return post, nil, user.Balance, nil
	}

	adons, err := s.addonsDbRepo.Get(post.ID)
	if err != nil {
		return nil, nil, 0, err
	}

	return post, adons, user.Balance, nil
}

func (s postService) GetPostByUser(ctx *gin.Context, serviceName string) (*model.Post, *model.Adons, int, error) {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, nil, 0, err
	}

	return s.Get(user, serviceName)
}

func (s postService) UpdatePost(post *model.Post) error {
	return s.postDbRepo.Update(post)
}
