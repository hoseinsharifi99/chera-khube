package service

import (
	"chera_khube/internal/constant"
	"chera_khube/internal/dto"
	"chera_khube/internal/helper"
	"chera_khube/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OAuthService interface {
	LoginWithDivar(ctx *gin.Context, service string) string
	GetToken(code string, service string) (*dto.AccessTokenResponse, error)
	GetPhoneNumber(accessToken string, service string) (*dto.PhoneNumberResponse, error)
	AdsEntry(ctx *gin.Context, service string) string
}

type oAuthService struct {
	repository repository.OAuthRepository
	Config     *helper.ServiceConfig
	logger     *zap.Logger
}

func NewOAuthService(repository repository.OAuthRepository, Config *helper.ServiceConfig, logger *zap.Logger) OAuthService {
	return oAuthService{repository: repository, Config: Config, logger: logger}
}

func (s oAuthService) GetPhoneNumber(accessToken string, service string) (*dto.PhoneNumberResponse, error) {
	config := s.Config.GetDivarConfig(service)

	return s.repository.GetPhoneNumber(dto.PhoneNumber{
		BaseUrl:     config.OAuthPhoneNumber.BaseUrl,
		ApiKey:      config.ApiKey,
		AccessToken: accessToken,
	})
}

func (s oAuthService) LoginWithDivar(ctx *gin.Context, service string) string {
	postToken := ctx.Query("post_token")
	returnUrl := ctx.Query("return_url")
	if !helper.IsDivarLink(returnUrl) {
		returnUrl = "https://divar.ir/my-divar/my-posts"
	}
	state := postToken

	var url string
	if service == constant.ApartmentServiceName {
		url = fmt.Sprintf(s.Config.Divar.Apartment.OAuth.BaseUrl, s.Config.Divar.Apartment.OAuth.ResponseType, s.Config.Divar.Apartment.ClientID, s.Config.Divar.Apartment.RedirectUrl, s.Config.Divar.Apartment.Scopes+postToken, state)
	} else {
		url = fmt.Sprintf(s.Config.Divar.General.OAuth.BaseUrl, s.Config.Divar.General.OAuth.ResponseType, s.Config.Divar.General.ClientID, s.Config.Divar.General.RedirectUrl, s.Config.Divar.General.Scopes+postToken, state)
	}

	return url
}

func (s oAuthService) AdsEntry(ctx *gin.Context, service string) string {
	var url string
	config := s.Config.GetDivarConfig(service)

	if service == constant.ApartmentServiceName {
		url = fmt.Sprintf(config.OAuth.BaseUrl, s.Config.Yektanet.Apartment.ResponseType, s.Config.Yektanet.Apartment.ClientID, s.Config.Yektanet.Apartment.RedirectUrl, "USER_POSTS_GET", uuid.New().String()+"__"+service)
	}

	return url
}

func (s oAuthService) GetToken(code string, service string) (*dto.AccessTokenResponse, error) {
	config := s.Config.GetDivarConfig(service)

	return s.repository.GetToken(dto.OAuthToken{
		BaseUrl:      config.OAuthToken.BaseUrl,
		Code:         code,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		GrantType:    config.OAuthToken.GrantType,
		RedirectUri:  config.RedirectUrl,
	})
}
