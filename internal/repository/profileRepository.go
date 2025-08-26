package repository

import "chera_khube/internal/model"

type ProfileRepository interface {
	Set(profile *model.Profile) (*model.Profile, error)
	Get(userID string) (*model.Profile, error)
	GetByUserID(userID uint) (*model.Profile, error)
	GetByPhoneNumber(phoneNumber string) (*model.Profile, error)
	Update(profile *model.Profile) error
	UpdateIsConnected(postToken string, isConnected bool) error
}
