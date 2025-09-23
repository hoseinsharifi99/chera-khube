package db

import (
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewPricingDB(db *gorm.DB, logger *zap.Logger) repository.PricingRepository {
	return &pricingDB{
		db:     db,
		logger: logger,
	}
}

type pricingDB struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (p pricingDB) ListWithFirstDiscount(serviceID uint) ([]*model.PricingLogic, error) {
	var pricing []*model.PricingLogic

	err := p.db.Where("id != ? AND id != ? AND service_id = ?", 1, 5, serviceID).Order("price ASC").Find(&pricing).Error
	if err != nil {
		return nil, err
	}

	return pricing, nil
}

func (p pricingDB) ListWithOutFirstDiscount(serviceID uint) ([]*model.PricingLogic, error) {
	var pricing []*model.PricingLogic

	err := p.db.Where("id != ? AND id != ? AND service_id = ?", 4, 6, serviceID).Order("price ASC").Find(&pricing).Error
	if err != nil {
		return nil, err
	}

	return pricing, nil
}

func (p pricingDB) Get(id int) (*model.PricingLogic, error) {
	var price model.PricingLogic
	err := p.db.Where("id = ?", id).Find(&price).Error
	if err != nil {
		return nil, err
	}

	return &price, nil
}

func (p pricingDB) FindByPrice(priceAmount uint) (*model.PricingLogic, error) {
	var price model.PricingLogic
	err := p.db.Where("price = ?", priceAmount).Find(&price).Error
	if err != nil {
		return nil, err
	}

	return &price, nil
}
