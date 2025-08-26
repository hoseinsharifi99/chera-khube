package router

import (
	"chera_khube/handler/controller"
	middlewares "chera_khube/handler/middleware"
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
)

type postRouter struct {
	postController controller.PostController
}

func NewPostRouter(postController controller.PostController) Router {
	return &postRouter{postController: postController}
}

func (r postRouter) HandleRoutes(router *gin.Engine, config *helper.ServiceConfig) {
	post := router.Group("v1").Group("post")
	post.GET("", middlewares.Jwt(config), r.postController.Get)
	post.GET("all", middlewares.Jwt(config), r.postController.GetAll)
}
