package repository

import (
	"chera_khube/internal/model"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Register(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	GetUserBalanceByPhoneNumber(phoneNumber string) (int, error)
	GetUserByPhoneNumberAndPassword(phoneNumber, password string) (*model.User, error)
	GetUserByID(userID uint) (*model.User, error)
}
