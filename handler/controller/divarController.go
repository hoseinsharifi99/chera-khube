package controller

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type DivarController interface {
	EditPost(ctx *gin.Context)
	EditAllPost(ctx *gin.Context)
	CreateAddons(ctx *gin.Context)
	GetAllDescription(ctx *gin.Context)
}
type divarController struct {
	divarService service.AddonsService
	config       *helper.ServiceConfig
	logger       *zap.Logger
}

func NewDivarController(
	divarService service.AddonsService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) DivarController {
	return &divarController{
		divarService: divarService,
		config:       config,
		logger:       logger,
	}
}

func (c divarController) EditPost(ctx *gin.Context) {
	postToken := ctx.Param("post-token")

	if postToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty post token"})
	}

	post, balance, err := c.divarService.AddWidgetToPost(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": post, "balance": balance})
}

func (c divarController) EditAllPost(ctx *gin.Context) {
	postToken := ctx.Param("post-token")

	if postToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty post token"})
	}

	post, balance, err := c.divarService.EditAllDescription(ctx, postToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": post, "balance": balance})
}

func (c divarController) CreateAddons(ctx *gin.Context) {
	//postToken := ctx.Param("post-token")
	//codes := ctx.Param("codes")

	//if postToken == "" {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty post token"})
	//}
	//
	//if codes == "" {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty codes"})
	//}
	//
	//post, balance, err := c.divarService.CreateAddons(ctx, postToken, codes)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	ctx.Abort()
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": "post", "balance": "balance"})
}

func (c divarController) GetAllDescription(ctx *gin.Context) {
	postToken := ctx.Param("post-token")

	if postToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty post token"})
	}

	post, balance, err := c.divarService.GetAllNewDesc(ctx, postToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "post": post, "balance": balance})
}
