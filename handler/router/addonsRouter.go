package router

import (
	"chera_khube/handler/controller"
	middlewares "chera_khube/handler/middleware"
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
)

type addonsRouter struct {
	addonsController controller.AddonsController
}

func NewAddonsRouter(addonsController controller.AddonsController) Router {
	return &addonsRouter{addonsController: addonsController}
}

func (r addonsRouter) HandleRoutes(router *gin.Engine, config *helper.ServiceConfig) {
	user := router.Group("v1").Group("addons")
	user.GET("create/:codes/:service", middlewares.Jwt(config), r.addonsController.CreateAddons)
	user.POST("add/widget/:service", middlewares.Jwt(config), r.addonsController.AddAddons)
	user.GET("widget/:service", middlewares.Jwt(config), r.addonsController.GetAddons)
	user.DELETE("widget/:service", middlewares.Jwt(config), r.addonsController.DeleteWidget)
	user.GET("config/:service", r.addonsController.GetConfig)
}
