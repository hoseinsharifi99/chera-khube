package router

import (
	"chera_khube/handler/controller"
	middlewares "chera_khube/handler/middleware"
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
)

type pricingRouter struct {
	pricingController controller.PricingController
}

func NewPricingRouter(pricingController controller.PricingController) Router {
	return &pricingRouter{pricingController: pricingController}
}

func (r pricingRouter) HandleRoutes(router *gin.Engine, config *helper.ServiceConfig) {
	post := router.Group("v1").Group("pricing")
	post.GET("", middlewares.Jwt(config), r.pricingController.List)
}
