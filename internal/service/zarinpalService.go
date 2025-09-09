package service

import (
	"chera_khube/internal/dto"
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type ZarinpalService interface {
	NewPaymentRequest(ctx *gin.Context, amount int, serviceName string) (*dto.PaymentResponseDto, error)
	PaymentVerification(ctx *gin.Context, d dto.PaymentVerificationDto) (*dto.PaymentVerificationResponseDto, error)
}

type zarinpalService struct {
	api                repository.ZarinpalRepository
	userPaymentService UserPaymentService
	userService        UserService
	pricingService     PricingService
	config             *helper.ServiceConfig
	logger             *zap.Logger
}

func NewZarinpalService(
	api repository.ZarinpalRepository,
	userPaymentService UserPaymentService,
	userService UserService,
	pricingService PricingService,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) ZarinpalService {
	return &zarinpalService{
		api:                api,
		userPaymentService: userPaymentService,
		userService:        userService,
		pricingService:     pricingService,
		config:             config,
		logger:             logger,
	}
}

func (s *zarinpalService) NewPaymentRequest(ctx *gin.Context, planID int, serviceName string) (*dto.PaymentResponseDto, error) {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, err
	}

	price, err := s.pricingService.Get(planID)
	if err != nil {
		return nil, err
	}

	d := dto.PaymentRequestDto{
		Amount:      int(price.Price) * 10,
		CallbackUrl: fmt.Sprintf(s.config.Zarinpal.CallbackUrl, user.ID, int(price.Price*10), planID, serviceName),
		Description: "اگهی پلاس",
		Email:       fmt.Sprintf("%d@user.com", user.ID),
		PhoneNumber: user.PhoneNumber,
	}

	if d.Amount < 1 {
		err = errors.New("amount must be a positive number")
		return nil, err
	}
	if d.CallbackUrl == "" {
		err = errors.New("callbackURL should not be empty")
		return nil, err
	}
	if d.PhoneNumber == "" {
		err = errors.New("phoneNumber should not be empty")
		return nil, err
	}

	return s.api.NewPaymentRequest(d)
}

func (s *zarinpalService) PaymentVerification(ctx *gin.Context, d dto.PaymentVerificationDto) (*dto.PaymentVerificationResponseDto, error) {
	log.Println(ctx.Request.URL)
	userIDString := ctx.Query("user_id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return nil, err
	}
	planIdString := ctx.Query("id")
	planID, err := strconv.Atoi(planIdString)
	if err != nil {
		return nil, err
	}
	serviceName := ctx.Query("service")

	user, err := s.userService.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	if d.Amount <= 0 {
		err = errors.New("amount must be a positive number")
		return nil, err
	}

	if d.Authority == "" {
		err = errors.New("authority should not be empty")
		return nil, err
	}

	price, err := s.pricingService.Get(planID)
	if err != nil {
		return nil, err
	}

	if d.Amount/10 != int(price.Price) {
		return nil, errors.New("amount is incorrect")
	}

	result, err := s.api.PaymentVerification(d)
	if err != nil {
		return nil, err
	}

	//Update user balance

	user.Balance = user.Balance + int(price.Token)
	user.IsUse = true

	user, err = s.userService.Update(user)
	if err != nil {
		return nil, err
	}

	status := result.StatusCode == 100

	up := &model.UserPayment{
		UserID: user.ID,
		Amount: d.Amount,
		RefID:  result.RefID,
		Status: status,
	}
	err = s.userPaymentService.Create(up)
	if err != nil {
		return nil, err
	}

	result.ServiceName = serviceName

	return result, nil
}
