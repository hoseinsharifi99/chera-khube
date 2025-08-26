package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewPostDb(db *gorm.DB, logger *zap.Logger) repository.PostDbRepo {
	return &postDb{
		db:     db,
		logger: logger,
	}
}

type postDb struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r postDb) Insert(post *model.Post) (*model.Post, error) {
	err := r.db.FirstOrCreate(post, "token = ?", post.Token).Error

	return post, err
}

func (r postDb) Update(post *model.Post) error {
	err := r.db.Save(post).Error

	return err
}

func (r postDb) Get(token string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("token = ?", token).First(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}
