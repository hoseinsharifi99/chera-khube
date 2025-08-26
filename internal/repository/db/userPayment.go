package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewUserPaymentRepository(db *gorm.DB, logger *zap.Logger) repository.UserPaymentRepository {
	return &userPaymentRepository{
		db:     db,
		logger: logger,
	}
}

type userPaymentRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r userPaymentRepository) Create(payment *model.UserPayment) error {
	return r.db.Create(payment).Error
}

func (r userPaymentRepository) List(userID uint) ([]*model.UserPayment, error) {
	var payments []*model.UserPayment
	err := r.db.Where("user_id = ?", userID).Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
