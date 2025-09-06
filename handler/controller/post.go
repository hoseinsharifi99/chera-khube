package controller

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PostController interface {
	Get(ctx *gin.Context)
}
type postController struct {
	postService service.PostService
	config      *helper.ServiceConfig
	logger      *zap.Logger
}

func (c postController) Get(ctx *gin.Context) {
	//post, addons, balance, err := c.postService.GetPostByUser(ctx)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	ctx.Abort()
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func NewPostController(
	postService service.PostService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) PostController {
	return &postController{
		postService: postService,
		config:      config,
		logger:      logger,
	}
}
