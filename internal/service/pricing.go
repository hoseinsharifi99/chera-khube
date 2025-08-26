package service

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PricingService interface {
	List(ctx *gin.Context, serviceID uint) ([]*model.PricingLogic, error)
	Get(id int) (*model.PricingLogic, error)
	FindByPrice(price uint) (*model.PricingLogic, error)
}

type pricingService struct {
	pricingRepo repository.PricingRepository
	userService UserService
	config      *helper.ServiceConfig
	logger      *zap.Logger
}

func NewPricingService(
	pricingRepo repository.PricingRepository,
	userService UserService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) PricingService {
	return &pricingService{
		pricingRepo: pricingRepo,
		userService: userService,
		config:      config,
		logger:      logger,
	}
}

func (p pricingService) List(ctx *gin.Context, serviceID uint) ([]*model.PricingLogic, error) {
	user, err := p.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, err
	}

	if user.IsUse == true {
		return p.pricingRepo.ListWithOutFirstDiscount(serviceID)
	}

	return p.pricingRepo.ListWithFirstDiscount(serviceID)
}

func (p pricingService) Get(id int) (*model.PricingLogic, error) {
	return p.pricingRepo.Get(id)
}

func (p pricingService) FindByPrice(price uint) (*model.PricingLogic, error) {
	return p.pricingRepo.FindByPrice(price)
}
