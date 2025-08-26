package controller

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PricingController interface {
	List(ctx *gin.Context)
}
type pricingController struct {
	pricingService service.PricingService
	config         *helper.ServiceConfig
	logger         *zap.Logger
}

func (c pricingController) List(ctx *gin.Context) {
	serviceString := ctx.Query("service_id")
	serviceID, _ := strconv.Atoi(serviceString)

	post, err := c.pricingService.List(ctx, uint(serviceID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "data": post})
}

func NewPricingController(
	pricingService service.PricingService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) PricingController {
	return &pricingController{
		pricingService: pricingService,
		config:         config,
		logger:         logger,
	}
}
