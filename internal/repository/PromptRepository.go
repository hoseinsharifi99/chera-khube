package repository

import "github.com/gin-gonic/gin"

type PromptRepository interface {
	Get(ctx *gin.Context, input string) (string, error)
}
