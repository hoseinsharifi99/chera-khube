package controller

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AddonsController interface {
	CreateAddons(ctx *gin.Context)
	AddAddons(ctx *gin.Context)
}
type addonsController struct {
	addonsService service.AddonsService
	config        *helper.ServiceConfig
	logger        *zap.Logger
}

func (a addonsController) CreateAddons(ctx *gin.Context) {
	codes := ctx.Param("codes")
	post, addons, balance, err := a.addonsService.CreateAddons(ctx, "postToken", codes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": post, "addons:": addons, "balance": balance})
}

func (a addonsController) AddAddons(ctx *gin.Context) {
	addons, balance, err := a.addonsService.AddWidgetToPost(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "addons": addons, "balance": balance})
}

func NewAddonsController(
	addonsService service.AddonsService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) AddonsController {
	return &addonsController{
		addonsService: addonsService,
		config:        config,
		logger:        logger,
	}
}
