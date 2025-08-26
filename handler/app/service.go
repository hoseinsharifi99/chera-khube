package app

import (
	bService "chera_khube/internal/service"
	"go.uber.org/zap"
)

type service struct {
	oauthService       bService.OAuthService
	userService        bService.UserService
	zarinpalService    bService.ZarinpalService
	userPaymentService bService.UserPaymentService
	postService        bService.PostService
	pricingService     bService.PricingService
	divarService       bService.AddonsService
	promptService      bService.PromptService
	addonsService      bService.AddonsService
}

func (a *application) InitService(repo *repository, logger *zap.Logger) *service {
	var srv service
	srv.oauthService = bService.NewOAuthService(repo.oauthRepository, a.config, logger)
	srv.userPaymentService = bService.NewUserPaymentService(repo.userPaymentRepository, logger)
	srv.userService = bService.NewUserService(srv.oauthService, repo.userRepository, srv.userPaymentService, repo.divarRepository, a.config, logger)
	srv.pricingService = bService.NewPricingService(repo.pricingRepository, srv.userService, a.config, logger)
	srv.zarinpalService = bService.NewZarinpalService(repo.zarinpalApiRepository, srv.userPaymentService, srv.userService, srv.pricingService, a.config, logger)
	srv.postService = bService.NewPostService(repo.postApiRepository, repo.postDBRepository, repo.addonsDbRepository, srv.userService, a.config, logger)
	srv.promptService = bService.NewPromptService(repo.promptRepository, logger)
	srv.divarService = bService.NewAddonsService(repo.divarRepository, srv.promptService, srv.userService, srv.postService, repo.addonsDbRepository, repo.widgetRepository, repo.configDbRepository, a.config, logger)
	srv.addonsService = bService.NewAddonsService(repo.divarRepository, srv.promptService, srv.userService, srv.postService, repo.addonsDbRepository, repo.widgetRepository, repo.configDbRepository, a.config, logger)
	return &srv
}
