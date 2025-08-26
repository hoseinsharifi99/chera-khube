package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB, logger *zap.Logger) repository.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r userRepository) Register(user *model.User) (*model.User, error) {
	err := r.db.FirstOrCreate(user, "phone_number = ?", user.PhoneNumber).Error

	return user, err
}

func (r userRepository) Update(user *model.User) (*model.User, error) {
	err := r.db.Save(user).Error

	return user, err
}

func (r userRepository) GetUserBalanceByPhoneNumber(phoneNumber string) (int, error) {
	var user model.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return 0, err
	}

	return user.Balance, nil
}

func (r userRepository) GetUserByPhoneNumberAndPassword(phoneNumber, password string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone_number = ? AND password = ?", phoneNumber, password).First(&user).Error; err != nil {
		return nil, repository.ErrUserNotFound
	}

	return &user, nil
}

func (r userRepository) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error

	return &user, err
}
