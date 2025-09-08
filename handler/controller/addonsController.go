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
	GetConfig(ctx *gin.Context)
	GetAddons(ctx *gin.Context)
	DeleteWidget(ctx *gin.Context)
}
type addonsController struct {
	addonsService service.AddonsService
	config        *helper.ServiceConfig
	logger        *zap.Logger
}

func (a addonsController) CreateAddons(ctx *gin.Context) {
	codes := ctx.Param("codes")
	srv := ctx.Param("service")

	post, addons, balance, err := a.addonsService.CreateAddons(ctx, "postToken", codes, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": post, "addons": addons, "balance": balance})
}

func (a addonsController) GetAddons(ctx *gin.Context) {
	srv := ctx.Param("service")
	addons, balance, err := a.addonsService.GetAddons(ctx, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "addons": addons, "balance": balance})
}

func (a addonsController) DeleteWidget(ctx *gin.Context) {
	srv := ctx.Param("service")

	err := a.addonsService.DeleteWidget(ctx, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a addonsController) AddAddons(ctx *gin.Context) {
	srv := ctx.Param("service")

	addons, balance, err := a.addonsService.AddWidgetToPost(ctx, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "addons": addons, "balance": balance})
}

func (a addonsController) GetConfig(ctx *gin.Context) {
	configs := a.addonsService.GetConfig()
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "configs": configs})
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
