package repository

import (
	"chera_khube/internal/model"
)

type PostApiRepo interface {
	Get(token string, erviceName string) (*model.Post, error)
}

type PostDbRepo interface {
	Insert(*model.Post) (*model.Post, error)
	Get(token string) (*model.Post, error)
	Update(post *model.Post) error
}
