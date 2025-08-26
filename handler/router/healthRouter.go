package router

import (
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthRouter struct {
}

func NewHealthRouter() Router {
	return &healthRouter{}
}

func (r healthRouter) HandleRoutes(router *gin.Engine, config *helper.ServiceConfig) {
	health := router.Group("")
	health.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "Service is working ..."}) })
	health.HEAD("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "Service is working ..."}) })
}
