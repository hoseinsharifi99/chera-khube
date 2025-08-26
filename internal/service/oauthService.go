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
	GetToken(code string) (*dto.AccessTokenResponse, error)
	GetTokenWithCustomRedirectUrl(code string, clientID string, clientSecret string, redirectUrl string) (*dto.AccessTokenResponse, error)
	GetProfileToken(code string) (*dto.AccessTokenResponse, error)
	GetAgahiToken(code string) (*dto.AccessTokenResponse, error)
	GetPhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error)
	GetProfilePhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error)
	GetAgahiPhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error)
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

func (s oAuthService) GetToken(code string) (*dto.AccessTokenResponse, error) {
	return s.repository.GetToken(dto.OAuthToken{
		BaseUrl:      s.Config.Divar.OAuthToken.BaseUrl,
		Code:         code,
		ClientID:     s.Config.Divar.ClientID,
		ClientSecret: s.Config.Divar.ClientSecret,
		GrantType:    s.Config.Divar.OAuthToken.GrantType,
		RedirectUri:  s.Config.Divar.RedirectUrl,
	})
}

func (s oAuthService) GetAgahiToken(code string) (*dto.AccessTokenResponse, error) {
	return s.repository.GetToken(dto.OAuthToken{
		BaseUrl:      s.Config.AgahiPlus.OAuthToken.BaseUrl,
		Code:         code,
		ClientID:     s.Config.AgahiPlus.ClientID,
		ClientSecret: s.Config.AgahiPlus.ClientSecret,
		GrantType:    s.Config.AgahiPlus.OAuthToken.GrantType,
		RedirectUri:  s.Config.AgahiPlus.RedirectUrl,
	})
}

func (s oAuthService) GetProfileToken(code string) (*dto.AccessTokenResponse, error) {
	return s.repository.GetToken(dto.OAuthToken{
		BaseUrl:      s.Config.CarDivar.OAuthToken.BaseUrl,
		Code:         code,
		ClientID:     s.Config.CarDivar.ClientID,
		ClientSecret: s.Config.CarDivar.ClientSecret,
		GrantType:    s.Config.CarDivar.OAuthToken.GrantType,
		RedirectUri:  s.Config.CarDivar.RedirectUrl,
	})
}

func (s oAuthService) GetPhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error) {
	return s.repository.GetPhoneNumber(dto.PhoneNumber{
		BaseUrl:     s.Config.Divar.OAuthPhoneNumber.BaseUrl,
		ApiKey:      s.Config.Divar.ApiKey,
		AccessToken: accessToken,
	})
}

func (s oAuthService) GetProfilePhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error) {
	return s.repository.GetPhoneNumber(dto.PhoneNumber{
		BaseUrl:     s.Config.CarDivar.OAuthPhoneNumber.BaseUrl,
		ApiKey:      s.Config.CarDivar.ApiKey,
		AccessToken: accessToken,
	})
}

func (s oAuthService) GetAgahiPhoneNumber(accessToken string) (*dto.PhoneNumberResponse, error) {
	return s.repository.GetPhoneNumber(dto.PhoneNumber{
		BaseUrl:     s.Config.AgahiPlus.OAuthPhoneNumber.BaseUrl,
		ApiKey:      s.Config.AgahiPlus.ApiKey,
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
	if service == constant.AgahiPlusServiceName {
		url = fmt.Sprintf(s.Config.AgahiPlus.OAuth.BaseUrl, s.Config.AgahiPlus.OAuth.ResponseType, s.Config.AgahiPlus.ClientID, s.Config.AgahiPlus.RedirectUrl, "USER_PHONE POST_EDIT."+postToken, state)
	} else if service == constant.LinkPlusServiceName {
		url = fmt.Sprintf(s.Config.CarDivar.OAuth.BaseUrl, s.Config.CarDivar.OAuth.ResponseType, s.Config.CarDivar.ClientID, s.Config.CarDivar.RedirectUrl, "USER_PHONE offline_access USER_ADDON_CREATE", state)
	} else {
		url = fmt.Sprintf(s.Config.Divar.OAuth.BaseUrl, s.Config.Divar.OAuth.ResponseType, s.Config.Divar.ClientID, s.Config.Divar.RedirectUrl, "USER_PHONE POST_EDIT."+postToken, state)
	}

	return url
}

func (s oAuthService) AdsEntry(ctx *gin.Context, service string) string {
	var url string
	if service == constant.AgahiPlusServiceName {
		url = fmt.Sprintf(s.Config.Divar.OAuth.BaseUrl, s.Config.Yektanet.AgahiPlus.ResponseType, s.Config.Yektanet.AgahiPlus.ClientID, s.Config.Yektanet.AgahiPlus.RedirectUrl, "USER_POSTS_GET", uuid.New().String()+"__"+service)
	} else if service == constant.LinkPlusServiceName {
		url = fmt.Sprintf(s.Config.Divar.OAuth.BaseUrl, s.Config.Yektanet.LinkPlus.ResponseType, s.Config.Yektanet.LinkPlus.ClientID, s.Config.Yektanet.LinkPlus.RedirectUrl, "USER_POSTS_GET", uuid.New().String()+"__"+service)
	} else {
		url = fmt.Sprintf(s.Config.Divar.OAuth.BaseUrl, s.Config.Yektanet.Apartment.ResponseType, s.Config.Yektanet.Apartment.ClientID, s.Config.Yektanet.Apartment.RedirectUrl, "USER_POSTS_GET", uuid.New().String()+"__"+service)
	}

	return url
}

func (s oAuthService) GetTokenWithCustomRedirectUrl(code string, clientID string, clientSecret string, redirectUrl string) (*dto.AccessTokenResponse, error) {
	return s.repository.GetToken(dto.OAuthToken{
		BaseUrl:      s.Config.Divar.OAuthToken.BaseUrl,
		Code:         code,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    s.Config.Divar.OAuthToken.GrantType,
		RedirectUri:  redirectUrl,
	})
}
