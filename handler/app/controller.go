package app

import (
	bController "chera_khube/handler/controller"
	"go.uber.org/zap"
)

type controller struct {
	userController    bController.UserController
	paymentController bController.PaymentController
	postController    bController.PostController
	pricingController bController.PricingController
	divarController   bController.DivarController
}

func (a *application) InitController(srv *service, repo *repository, logger *zap.Logger) *controller {
	var ctrl controller
	ctrl.userController = bController.NewUserController(srv.userService, repo.appLogRepository, a.config, logger)
	ctrl.paymentController = bController.NewPaymentController(srv.zarinpalService, a.config, logger)
	ctrl.postController = bController.NewPostController(srv.postService, a.config, logger)
	ctrl.pricingController = bController.NewPricingController(srv.pricingService, a.config, logger)
	ctrl.divarController = bController.NewDivarController(srv.divarService, a.config, logger)

	return &ctrl
}
