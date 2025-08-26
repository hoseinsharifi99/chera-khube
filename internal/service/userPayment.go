package service

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"go.uber.org/zap"
)

type UserPaymentService interface {
	Create(*model.UserPayment) error
	List(userID uint) ([]*model.UserPayment, error)
}

type userPayment struct {
	userPaymentRepo repository.UserPaymentRepository
	logger          *zap.Logger
}

func NewUserPaymentService(userPaymentRepo repository.UserPaymentRepository, logger *zap.Logger) UserPaymentService {
	return &userPayment{
		userPaymentRepo: userPaymentRepo,
		logger:          logger,
	}
}

func (u userPayment) Create(payment *model.UserPayment) error {
	return u.userPaymentRepo.Create(payment)
}

func (u userPayment) List(userID uint) ([]*model.UserPayment, error) {
	return u.userPaymentRepo.List(userID)
}
