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
	user.GET("create/:codes", middlewares.Jwt(config), r.addonsController.CreateAddons)
	user.GET("widget", middlewares.Jwt(config), r.addonsController.AddAddons)
	user.GET("/config", r.addonsController.GetConfig)
}
