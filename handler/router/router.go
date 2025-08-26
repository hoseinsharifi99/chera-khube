package router

import (
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
)

type Router interface {
	HandleRoutes(router *gin.Engine, config *helper.ServiceConfig)
}
