package service

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"strings"
)

type AddonsService interface {
	AddWidgetToPost(ctx *gin.Context, serviceName string) (*model.Adons, float64, error)
	CreateAddons(ctx *gin.Context, postToken, codes string, serviceName string) (*model.Post, *model.Adons, int, error)
	GetAllNewDesc(ctx *gin.Context, postToken string) (*model.Post, int, error)
	AddWidget(postToken, accessToken string, wid map[string]string, addons model.Adons, serviceName string) error
	DeleteWidget(ctx *gin.Context, serviceName string) error
	GetConfig() []model.Config
	GetAddons(ctx *gin.Context) (*model.Adons, int, error)
}

type addonsService struct {
	divarRepo     repository.DivarRepository
	promptService PromptService
	userService   UserService
	postService   PostService
	widgetRepo    repository.WidgetRepository
	addonDbRepo   repository.AdonsDbRepo
	configRepo    repository.ConfigRepository
	config        *helper.ServiceConfig
	logger        *zap.Logger
}

func NewAddonsService(
	divarRepo repository.DivarRepository,
	promptService PromptService,
	userService UserService,
	postService PostService,
	addonDbRepo repository.AdonsDbRepo,
	widgetRepo repository.WidgetRepository,
	configRepo repository.ConfigRepository,
	config *helper.ServiceConfig,
	logger *zap.Logger,
) AddonsService {
	return &addonsService{
		divarRepo:     divarRepo,
		promptService: promptService,
		userService:   userService,
		postService:   postService,
		addonDbRepo:   addonDbRepo,
		widgetRepo:    widgetRepo,
		configRepo:    configRepo,
		config:        config,
		logger:        logger,
	}
}

func (s addonsService) AddWidgetToPost(ctx *gin.Context, serviceName string) (*model.Adons, float64, error) {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	post, ad, _, err := s.postService.GetPostByUser(ctx)
	if err != nil {
		return nil, 0, err
	}

	if ad.Codes == nil {
		return ad, float64(user.Balance), errors.New("code is nil")
	}
	keys := strings.Split(*ad.Codes, ",")

	AConf := s.configRepo.ListAsMap()
	ConfMap := make(map[string]string)
	for _, c := range keys {
		ConfMap[c] = AConf[c]
	}

	err = s.AddWidget(post.Token, user.AccessToken, ConfMap, *ad, serviceName)
	if err != nil {
		log.Println("error on add widget", err.Error())
		return ad, float64(user.Balance), err
	}

	ad.IsConnected = true

	err = s.addonDbRepo.Update(ad)
	if err != nil {
		log.Println("error on update db", err.Error())
	}

	return ad, float64(user.Balance), err
}

func (s addonsService) AddWidgetToPostAfterCreate(ad *model.Adons, post *model.Post, accessToken string, serviceName string) (*model.Adons, error) {
	if ad.Codes == nil {
		return ad, errors.New("code is nil")
	}
	keys := strings.Split(*ad.Codes, ",")

	AConf := s.configRepo.ListAsMap()
	ConfMap := make(map[string]string)
	for _, c := range keys {
		ConfMap[c] = AConf[c]
	}

	err := s.AddWidget(post.Token, accessToken, ConfMap, *ad, serviceName)
	if err != nil {
		log.Println("error on add widget", err.Error())
		return ad, err
	}

	ad.IsConnected = true

	err = s.addonDbRepo.Update(ad)
	if err != nil {
		log.Println("error on update db", err.Error())
	}

	return ad, err
}

func (s addonsService) GetConfig() []model.Config {
	return s.configRepo.List()
}

func (s addonsService) CreateAddons(ctx *gin.Context, postToken, codes string, serviceName string) (*model.Post, *model.Adons, int, error) {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, nil, 0, err
	}

	if user.Balance < 1 {
		return nil, nil, 0, errors.New("not enough balance")
	}

	post, _, balance, err := s.postService.Get(user)
	if err != nil {
		return nil, nil, 0, err
	}

	addonsInfo := s.configRepo.GetByCodes(strings.Split(codes, ","))
	addonsDesc := ""
	for _, a := range addonsInfo {
		addonsDesc += a.Description
		addonsDesc += "و"
	}

	newDescription, err := s.promptService.CreateNewDescription(ctx, post.Data, addonsDesc)
	if err != nil {
		return nil, nil, 0, err
	}

	addons := &model.Adons{
		PostID:      post.ID,
		IsConnected: false,
		Description: newDescription,
		Codes:       &codes,
	}

	_, err = s.addonDbRepo.Insert(addons)

	balance = user.Balance - 1
	err = s.userService.UpdateBalance(user, balance)
	if err != nil {
		return nil, nil, 0, err
	}

	addons, err = s.AddWidgetToPostAfterCreate(addons, post, user.AccessToken, serviceName)
	if err != nil {
		return nil, nil, user.Balance, err
	}

	return post, addons, balance, nil
}

func (s addonsService) GetAllNewDesc(ctx *gin.Context, postToken string) (*model.Post, int, error) {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	if user.Balance < 1 {
		return nil, 0, errors.New("not enough balance")
	}

	post, _, _, err := s.postService.Get(user)
	if err != nil {
		return nil, 0, err
	}

	_, err = s.promptService.CreateAgahiNewDescription(ctx, post.Data)
	if err != nil {
		return nil, 0, err
	}

	err = s.postService.UpdatePost(post)

	balance := user.Balance - 1
	err = s.userService.UpdateBalance(user, balance)
	if err != nil {
		return nil, 0, err
	}

	return post, balance, nil
}

func (s addonsService) AddWidget(postToken, accessToken string, wid map[string]string, addons model.Adons, serviceName string) error {
	widget := s.createWidgets(wid, addons)

	err := s.widgetRepo.Send(widget, postToken, accessToken, serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (s addonsService) DeleteWidget(ctx *gin.Context, serviceName string) error {
	user, err := s.userService.GetUserWithContext(ctx)
	if err != nil {
		return err
	}

	err = s.widgetRepo.Delete(user.PostToken, user.AccessToken, serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (s addonsService) DeleteByPostTokenWidget(postToken, accessToken string, serviceName string) error {
	err := s.widgetRepo.Delete(postToken, accessToken, serviceName)
	if err != nil {
		return err
	}

	return nil
}

func (s addonsService) createWidgets(wd map[string]string, addons model.Adons) *model.DivarWidget {
	var dWidget model.DivarWidget
	for key, val := range wd {
		var widget model.EventWidget
		widget.EventRow.Title = val
		widget.EventRow.IconName = key
		widget.EventRow.HasDivider = false

		dWidget.Widgets = append(dWidget.Widgets, widget)
	}

	var descriptionWidget model.DescriptionRow
	descriptionWidget.Text = addons.Description
	descriptionWidget.HasDivider = false
	descriptionWidget.Expandable = false

	dWidget.Widgets = append(dWidget.Widgets, descriptionWidget)

	return &dWidget
}

func (s addonsService) GetAddons(ctx *gin.Context) (*model.Adons, int, error) {
	_, addons, balance, err := s.postService.GetPostByUser(ctx)
	if err != nil {
		return nil, 0, err
	}

	return addons, balance, nil
}
