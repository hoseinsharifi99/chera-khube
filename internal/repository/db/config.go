package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewConfigDb(db *gorm.DB, logger *zap.Logger) repository.ConfigRepository {
	return &configDb{
		db:     db,
		logger: logger,
	}
}

type configDb struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (r configDb) List(serviceName string) []model.Config {
	var configs []model.Config
	err := r.db.Where("service_name = ?", serviceName).Find(&configs).Error
	if err != nil {
		r.logger.Error("failed to fetch configs", zap.Error(err))
		return []model.Config{}
	}
	return configs
}

func (r configDb) ListAsMap(serviceName string) map[string]string {
	var configs []model.Config
	result := make(map[string]string)

	err := r.db.Where("service_name = ?", serviceName).Find(&configs).Error
	if err != nil {
		r.logger.Error("failed to fetch configs", zap.Error(err))
		return result
	}

	for _, cfg := range configs {
		result[cfg.Code] = cfg.Description
	}

	return result
}

func (r configDb) GetByCodes(codes []string, serviceName string) []model.Config {
	var configs []model.Config
	err := r.db.Where("service_name = ? AND code IN ?", serviceName, codes).Find(&configs).Error
	if err != nil {
		r.logger.Error("failed to fetch configs by codes", zap.Error(err), zap.Strings("codes", codes))
		return []model.Config{}
	}
	return configs
}
