package controller

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/repository"
	"chera_khube/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type UserController interface {
	LoginWithDivar(ctx *gin.Context)
	OAuth(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetBalance(ctx *gin.Context)
	HasBalance(ctx *gin.Context)
	CheckLogin(ctx *gin.Context)
	CallOAuth(ctx *gin.Context)
	GetPosts(ctx *gin.Context)
	AdsEntry(ctx *gin.Context)
	AdsOAuth(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	appLog      repository.AppLogRepository
	config      *helper.ServiceConfig
	logger      *zap.Logger
}

func NewUserController(userService service.UserService, appLog repository.AppLogRepository, config *helper.ServiceConfig, logger *zap.Logger) UserController {
	return &userController{userService: userService, appLog: appLog, config: config, logger: logger}
}

func (u userController) LoginWithDivar(ctx *gin.Context) {
	serviceName := ctx.Param("service")
	postToken := ctx.Query("post_token")
	returnUrl := ctx.Query("return_url")

	redirectUrl := fmt.Sprintf(u.config.App.FrontEndEntryRedirect+"/"+serviceName+"?post_token=%s&return_url=%s&service=%s", postToken, returnUrl, serviceName)
	fmt.Println("redirect to:", returnUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	return
}

func (u userController) CallOAuth(ctx *gin.Context) {
	serviceName := ctx.Param("service")
	redirectUrl := u.userService.LoginWithDivar(ctx, serviceName)
	log.Println("call oauth:", redirectUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (u userController) OAuth(ctx *gin.Context) {
	srv := ctx.Param("service")

	user, err := u.userService.OAuth(ctx, srv)
	postToken := ctx.Query("state")

	if err != nil {
		log.Println(err.Error())
		//redirect to access denied page:
		redirectUrl := fmt.Sprintf(u.config.App.FrontEndAccessDeniedRedirect, postToken, srv)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		return
	}

	redirectUrl := fmt.Sprintf(u.config.App.FrontEndLoginRedirect, srv, user.PhoneNumber, user.JwtToken, user.JwtExpireAt.Format("2006-01-02 15:04:05"), user.PostToken)
	ctx.SetCookie("token", user.JwtToken, int(user.JwtExpireAt.Sub(time.Now()).Seconds()), "/", "", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (u userController) Register(ctx *gin.Context) {
	if !u.config.System.DevelopMode {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "develop mode is not enabled"})
		ctx.Abort()
		return
	}
	user, err := u.userService.Register()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
}
func (u userController) Login(ctx *gin.Context) {
	if !u.config.System.DevelopMode {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "develop mode is not enabled"})
		ctx.Abort()
		return
	}
	phoneNumber := ctx.Query("phone_number")
	password := ctx.Query("password")

	user, err := u.userService.Login(phoneNumber, password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
}

func (u userController) GetBalance(ctx *gin.Context) {
	balanceAndUserPayments, err := u.userService.GetUserBalance(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "info": balanceAndUserPayments})
}

func (u userController) HasBalance(ctx *gin.Context) {
	hasBalance, err := u.userService.HasBalance(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "has_balance": hasBalance})
}

func (u userController) CheckLogin(ctx *gin.Context) {
	_, err := u.userService.GetUserWithContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "result": false})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "result": true})
}

func (u userController) GetPosts(ctx *gin.Context) {
	accessToken := ctx.Query("access_token")
	srv := ctx.Param("service")

	resp, err := u.userService.GetPosts(accessToken, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (u userController) AdsEntry(ctx *gin.Context) {
	srv := ctx.Param("service")
	redirectUrl := u.userService.AdsEntry(ctx, srv)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (u userController) AdsOAuth(ctx *gin.Context) {
	srv := ctx.Param("service")
	accessToken, err := u.userService.AdsOAuth(ctx, srv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(u.config.Yektanet.FrontRedirectUrl, accessToken, srv))
}
