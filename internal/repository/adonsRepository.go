package repository

import "chera_khube/internal/model"

type AdonsDbRepo interface {
	Insert(adons *model.Adons) (*model.Adons, error)
	Get(postID uint, serviceName string) (*model.Adons, error)
	Update(adons *model.Adons) error
	DeleteByID(id uint) error
}
