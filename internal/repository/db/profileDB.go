package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func NewProfileDB(db *gorm.DB, logger *zap.Logger) repository.ProfileRepository {
	return &profileDB{
		db:     db,
		logger: logger,
	}
}

type profileDB struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r profileDB) Set(profile *model.Profile) (*model.Profile, error) {
	err := r.db.Save(profile).Where("post_token = ?", profile.PostToken).Error

	return profile, err
}

func (r profileDB) Get(postToken string) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Where("post_token = ?", postToken).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r profileDB) GetByPhoneNumber(phoneNumber string) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Where("phone_number = ?", phoneNumber).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r profileDB) GetByUserID(userID uint) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r profileDB) Update(profile *model.Profile) error {
	return r.db.Debug().Save(profile).Error
}

func (r profileDB) UpdateIsConnected(postToken string, isConnected bool) error {
	result := r.db.Model(&model.Profile{}).
		Where("phone_number = ?", postToken).
		Update("is_connected", isConnected)

	if result.Error != nil {
		r.logger.Error("failed to update is_connected", zap.String("postToken", postToken), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("no profile found to update is_connected", zap.String("postToken", postToken))
		return gorm.ErrRecordNotFound
	}

	return nil
}

func getSMSProfiles(db *gorm.DB) ([]*model.Profile, error) {
	twoDaysAgo := time.Now().Add(-48 * time.Hour)

	var profiles []*model.Profile
	err := db.Where("connected_until < ? AND is_connected = ?", twoDaysAgo, true).Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	return profiles, nil
}
