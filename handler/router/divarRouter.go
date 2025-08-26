package router

import (
	"chera_khube/handler/controller"
	middlewares "chera_khube/handler/middleware"
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
)

type divarRouter struct {
	divarController controller.DivarController
}

func NewDivarRouter(divarController controller.DivarController) Router {
	return &divarRouter{divarController: divarController}
}

func (r divarRouter) HandleRoutes(router *gin.Engine, config *helper.ServiceConfig) {
	user := router.Group("v1").Group("divar")
	user.GET("edit/:post-token", middlewares.Jwt(config), r.divarController.EditPost)
	user.GET("edit/all/:post-token", middlewares.Jwt(config), r.divarController.EditAllPost)
	user.GET("description/:post-token", middlewares.Jwt(config), r.divarController.CreateAddons)
	user.GET("description/all/:post-token", middlewares.Jwt(config), r.divarController.GetAllDescription)
}
