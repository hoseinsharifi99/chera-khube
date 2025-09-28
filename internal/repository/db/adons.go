package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewAdonsDb(db *gorm.DB, logger *zap.Logger) repository.AdonsDbRepo {
	return &adonsDb{
		db:     db,
		logger: logger,
	}
}

type adonsDb struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r adonsDb) Insert(adons *model.Adons) (*model.Adons, error) {
	err := r.db.FirstOrCreate(adons, "post_id = ?", adons.PostID).Error

	return adons, err
}

func (r adonsDb) Update(adons *model.Adons) error {
	err := r.db.Save(adons).Error

	return err
}

func (r adonsDb) Get(postID uint, serviceName string) (*model.Adons, error) {
	var adons model.Adons
	err := r.db.Where("post_id = ? AND service = ?", postID, serviceName).First(&adons).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &adons, nil
}

func (r adonsDb) DeleteByID(id uint) error {
	err := r.db.Delete(&model.Adons{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // nothing to delete
		}
		return err
	}
	return nil
}
