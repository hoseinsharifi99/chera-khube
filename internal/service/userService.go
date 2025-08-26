package service

import (
	middlewares "chera_khube/handler/middleware"
	"chera_khube/internal/constant"
	"chera_khube/internal/dto"
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"chera_khube/internal/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"strings"
	"time"
)

type UserService interface {
	OAuth(ctx *gin.Context) (*model.User, error)
	ProfileOAuth(ctx *gin.Context) (*model.User, error)
	AgahiOAuth(ctx *gin.Context) (*model.User, error)
	CheckEnoughBalance(phoneNumber string) (enoughBalance bool, err error)
	GetUserByPhoneNumberAndPassword(phoneNumber, password string) (*model.User, error)
	Register() (*model.User, error)
	Login(phoneNumber, password string) (*model.User, error)
	UpdateBalance(user *model.User, newBalance int) error
	GetUserWithContext(ctx *gin.Context) (*model.User, error)
	GetUserBalance(ctx *gin.Context) (*response.UserBalancePayment, error)
	LoginWithDivar(ctx *gin.Context, service string) string
	HasBalance(ctx *gin.Context) (bool, error)
	GetUserByID(userID uint) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	AdsEntry(ctx *gin.Context, service string) string
	AdsOAuth(ctx *gin.Context, service string) (string, error)
	GetPosts(accessToken string, serviceName string) (*response.GetPostsResponse, error)
}

type userService struct {
	oAuthService       OAuthService
	userRepository     repository.UserRepository
	userPaymentService UserPaymentService
	divarApi           repository.DivarRepository
	config             *helper.ServiceConfig
	logger             *zap.Logger
}

func NewUserService(
	oAuthService OAuthService,
	userRepository repository.UserRepository,
	userPaymentService UserPaymentService,
	divarApi repository.DivarRepository,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) UserService {
	return &userService{
		oAuthService:       oAuthService,
		userRepository:     userRepository,
		userPaymentService: userPaymentService,
		divarApi:           divarApi,
		config:             config,
		logger:             logger,
	}
}

func (u userService) LoginWithDivar(ctx *gin.Context, service string) string {
	return u.oAuthService.LoginWithDivar(ctx, service)
}

func (u userService) OAuth(ctx *gin.Context) (*model.User, error) {
	fmt.Println(ctx.Request.URL)
	code := ctx.Query("code")

	if code == "" {
		return nil, errors.New("code is empty")
	}
	accessTokenResponse, err := u.oAuthService.GetToken(code)
	if err != nil {
		return nil, err
	}

	accessToken := accessTokenResponse.AccessToken
	//Register user to db

	phoneNumber, err := u.getPhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.register(*phoneNumber, accessTokenResponse, 1)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) ProfileOAuth(ctx *gin.Context) (*model.User, error) {
	fmt.Println(ctx.Request.URL)
	code := ctx.Query("code")

	if code == "" {
		return nil, errors.New("code is empty")
	}
	accessTokenResponse, err := u.oAuthService.GetProfileToken(code)
	if err != nil {
		return nil, err
	}

	accessToken := accessTokenResponse.AccessToken
	//Register user to db
	phoneNumber, err := u.getProfilePhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.profileRegister(*phoneNumber, accessTokenResponse)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) AgahiOAuth(ctx *gin.Context) (*model.User, error) {
	fmt.Println(ctx.Request.URL)
	code := ctx.Query("code")

	if code == "" {
		return nil, errors.New("code is empty")
	}
	accessTokenResponse, err := u.oAuthService.GetAgahiToken(code)
	if err != nil {
		return nil, err
	}

	accessToken := accessTokenResponse.AccessToken
	//Register user to db
	phoneNumber, err := u.getAgahiPhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.register(*phoneNumber, accessTokenResponse, 3)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) getPostToken(atResponse *dto.AccessTokenResponse) string {
	scopes := strings.Split(atResponse.Scope, " ")
	if len(scopes) > 2 {
		return "profile"
	}
	return strings.Split(scopes[1], ".")[1]
}
func (u userService) register(phoneNumber string, atResponse *dto.AccessTokenResponse, serviceID int) (*model.User, error) {
	if phoneNumber == "" {
		return nil, errors.New("phone number is empty")
	}

	pass := helper.GenerateTag()
	postToken := u.getPostToken(atResponse)
	user, err := u.userRepository.Register(&model.User{
		PhoneNumber: phoneNumber,
		Balance:     0,
	})

	if err != nil {
		return nil, err
	}

	user.PostToken = postToken
	user.AccessToken = atResponse.AccessToken
	user.ExpiresAt = (int64)(atResponse.ExpiresIn)
	user.Password = helper.GetPassword(pass, u.config.App.Salt)
	user.ServiceId = serviceID

	user, _ = u.userRepository.Update(user)

	token, expireAt := helper.GenerateAllToken(dto.Token{
		PhoneNumber: user.PhoneNumber,
		Password:    pass,
		SecretKey:   u.config.JWT.Secret,
		ExpireHour:  u.config.JWT.ExpireHour,
	})

	user.JwtToken = token
	user.JwtExpireAt = expireAt

	return user, nil
}

func (u userService) profileRegister(phoneNumber string, atResponse *dto.AccessTokenResponse) (*model.User, error) {
	if phoneNumber == "" {
		return nil, errors.New("phone number is empty")
	}

	pass := helper.GenerateTag()
	postToken := u.getPostToken(atResponse)
	user, err := u.userRepository.Register(&model.User{
		PhoneNumber: phoneNumber,
		Balance:     0,
	})

	if err != nil {
		return nil, err
	}

	user.PostToken = postToken
	user.RefreshToken = atResponse.RefreshToken
	user.AccessToken = atResponse.AccessToken
	user.ExpiresAt = (int64)(atResponse.ExpiresIn)
	user.ServiceId = 2
	user.Password = helper.GetPassword(pass, u.config.App.Salt)
	if user.IsFirstAgahi == false {
		user.ExpirationTime = time.Now().AddDate(0, 1, 0)
		user.IsFirstAgahi = true
	}

	user, _ = u.userRepository.Update(user)

	token, expireAt := helper.GenerateAllToken(dto.Token{
		PhoneNumber: user.PhoneNumber,
		Password:    pass,
		SecretKey:   u.config.JWT.Secret,
		ExpireHour:  u.config.JWT.ExpireHour,
	})

	user.JwtToken = token
	user.JwtExpireAt = expireAt

	return user, nil
}

func (u userService) getPhoneNumber(accessToken string) (phoneNumber *string, err error) {
	phoneNumbers, err := u.oAuthService.GetPhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	return &phoneNumbers.PhoneNumber, nil
}

func (u userService) getProfilePhoneNumber(accessToken string) (phoneNumber *string, err error) {
	phoneNumbers, err := u.oAuthService.GetProfilePhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	return &phoneNumbers.PhoneNumber, nil
}

func (u userService) getAgahiPhoneNumber(accessToken string) (phoneNumber *string, err error) {
	phoneNumbers, err := u.oAuthService.GetAgahiPhoneNumber(accessToken)
	if err != nil {
		return nil, err
	}

	return &phoneNumbers.PhoneNumber, nil
}

func (u userService) CheckEnoughBalance(phoneNumber string) (enoughBalance bool, err error) {
	balance, err := u.userRepository.GetUserBalanceByPhoneNumber(phoneNumber)
	if err != nil {
		return false, err
	}

	if balance < u.config.App.InquiryCost {
		return false, nil
	}

	return true, nil
}

func (u userService) GetUserByPhoneNumberAndPassword(phoneNumber, password string) (*model.User, error) {
	return u.userRepository.GetUserByPhoneNumberAndPassword(phoneNumber, helper.GetPassword(password, u.config.App.Salt))
}

func (u userService) Register() (*model.User, error) {
	token, expire := helper.GenerateAllToken(dto.Token{
		PhoneNumber: "test",
		Password:    "121",
		SecretKey:   u.config.JWT.Secret,
		ExpireHour:  u.config.JWT.ExpireHour,
	})
	return u.userRepository.Register(&model.User{
		PhoneNumber: "test",
		Balance:     10000000,
		AccessToken: "test",
		ExpiresAt:   100000000000000,
		Password:    helper.GetPassword("121", u.config.App.Salt),
		JwtToken:    token,
		JwtExpireAt: expire,
	})
}
func (u userService) Login(phoneNumber, password string) (*model.User, error) {
	user, err := u.userRepository.GetUserByPhoneNumberAndPassword(phoneNumber, helper.GetPassword("121", u.config.App.Salt))
	token, expire := helper.GenerateAllToken(dto.Token{
		PhoneNumber: "test",
		Password:    "121",
		SecretKey:   u.config.JWT.Secret,
		ExpireHour:  u.config.JWT.ExpireHour,
	})

	user.JwtToken = token
	user.JwtExpireAt = expire

	return user, err
}

func (u userService) UpdateBalance(user *model.User, newBalance int) error {
	if newBalance < 0 {
		return errors.New("not enough balance")
	}

	user.Balance = newBalance
	_, err := u.userRepository.Update(user)

	return err
}

func (u userService) GetUserWithContext(ctx *gin.Context) (*model.User, error) {
	phoneNumber := ctx.Keys[middlewares.PHONE_NUMBER].(string)
	password := ctx.Keys[middlewares.PASSWORD].(string)

	return u.GetUserByPhoneNumberAndPassword(phoneNumber, password)
}

func (u userService) GetUserBalance(ctx *gin.Context) (*response.UserBalancePayment, error) {
	user, err := u.GetUserWithContext(ctx)
	if err != nil {
		return nil, err
	}

	ups, err := u.userPaymentService.List(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.UserBalancePayment{
		Balance:      user.Balance,
		UserPayments: ups,
	}, nil
}

func (u userService) HasBalance(ctx *gin.Context) (bool, error) {
	user, err := u.GetUserWithContext(ctx)
	if err != nil {
		return false, err
	}

	if user.Balance < u.config.App.InquiryCost {
		return false, nil
	}

	return true, nil
}

func (u userService) GetUserByID(userID uint) (*model.User, error) {
	return u.userRepository.GetUserByID(userID)
}

func (u userService) Update(user *model.User) (*model.User, error) {
	return u.userRepository.Update(user)
}

func (u userService) AdsEntry(ctx *gin.Context, service string) string {
	return u.oAuthService.AdsEntry(ctx, service)
}

func (u userService) AdsOAuth(ctx *gin.Context, service string) (string, error) {
	code := ctx.Query("code")

	if code == "" {
		return "", errors.New("code is empty")
	}
	redirectUrl := ""
	clientID := ""
	clientSecret := ""
	switch service {
	case constant.AgahiPlusServiceName:
		redirectUrl = u.config.Yektanet.AgahiPlus.RedirectUrl
		clientID = u.config.Yektanet.AgahiPlus.ClientID
		clientSecret = u.config.AgahiPlus.ClientSecret
	case constant.LinkPlusServiceName:
		redirectUrl = u.config.Yektanet.LinkPlus.RedirectUrl
		clientID = u.config.Yektanet.LinkPlus.ClientID
		clientSecret = u.config.CarDivar.ClientSecret
	default:
		redirectUrl = u.config.Yektanet.Apartment.RedirectUrl
		clientID = u.config.Yektanet.Apartment.ClientID
		clientSecret = u.config.Divar.ClientSecret
	}

	log.Println("ADS OAUTH", service, redirectUrl)

	accessTokenResponse, err := u.oAuthService.GetTokenWithCustomRedirectUrl(code, clientID, clientSecret, redirectUrl)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	log.Println("CODE:", code, "ACCESS_TOKEN:", accessTokenResponse, accessTokenResponse.AccessToken)

	return accessTokenResponse.AccessToken, nil

}

func (u userService) GetPosts(accessToken string, serviceName string) (*response.GetPostsResponse, error) {
	apiKey := ""
	switch serviceName {
	case constant.AgahiPlusServiceName:
		apiKey = u.config.AgahiPlus.ApiKey
	case constant.LinkPlusServiceName:
		apiKey = u.config.CarDivar.ApiKey
	default:
		apiKey = u.config.Divar.ApiKey
	}
	return u.divarApi.GetPostTokens(u.config.Divar.GetPosts, apiKey, accessToken)
}
