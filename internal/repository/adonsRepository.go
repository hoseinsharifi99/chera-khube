package repository

import "chera_khube/internal/model"

type AdonsDbRepo interface {
	Insert(adons *model.Adons) (*model.Adons, error)
	Get(postID uint) (*model.Adons, error)
	Update(adons *model.Adons) error
}
