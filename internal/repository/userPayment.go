package repository

import "chera_khube/internal/model"

type UserPaymentRepository interface {
	Create(*model.UserPayment) error
	List(userID uint) ([]*model.UserPayment, error)
}
