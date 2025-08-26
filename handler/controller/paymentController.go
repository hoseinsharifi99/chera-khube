package controller

import (
	"chera_khube/internal/dto"
	"chera_khube/internal/helper"
	"chera_khube/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type PaymentController interface {
	NewPayment(ctx *gin.Context)
	Verify(ctx *gin.Context)
}

type paymentController struct {
	zarinpalService service.ZarinpalService
	config          *helper.ServiceConfig
	logger          *zap.Logger
}

func NewPaymentController(
	zarinpalService service.ZarinpalService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) PaymentController {
	return &paymentController{
		zarinpalService: zarinpalService,
		config:          config,
		logger:          logger,
	}
}

func (c paymentController) NewPayment(ctx *gin.Context) {
	var planID int
	var err error
	planIdString := ctx.Query("plan_id")
	if planIdString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("empty plan id")})
		ctx.Abort()
	}

	serviceName := ctx.Query("service")
	if serviceName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("empty service")})
		ctx.Abort()
	}

	planID, err = strconv.Atoi(planIdString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	result, err := c.zarinpalService.NewPaymentRequest(ctx, planID, serviceName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result.PaymentUrl})
}

func (c paymentController) Verify(ctx *gin.Context) {
	authority := ctx.Query("Authority")
	amountString := ctx.Query("amount")
	amount, err := strconv.Atoi(amountString)
	serviceName := ctx.Query("service")

	result, err := c.zarinpalService.PaymentVerification(ctx, dto.PaymentVerificationDto{
		Amount:    amount,
		Authority: authority,
	})

	if err != nil {

		redirectUrl := fmt.Sprintf(c.config.App.FrontEndPurchaseRedirect, 400, err.Error(), serviceName)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	redirectUrl := fmt.Sprintf(c.config.App.FrontEndPurchaseRedirect, result.StatusCode, result.RefID, serviceName)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
