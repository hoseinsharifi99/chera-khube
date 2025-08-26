package service

import (
	"chera_khube/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

var kamkardanchar = "این متن آگهی  بیشتر از ۱۰۰۰ کارکتر است. کمترین تعداد کارکتر ممکن رو حذف کن تا متن کمتر از ۱۰۰۰ کارکتر بشه. از حذف اطلاعات مهم آگهی خوداری کن."

type PromptService interface {
	CreateNewDescription(ctx *gin.Context, data, addons string) (string, error)
	CreateAgahiNewDescription(ctx *gin.Context, input string) (string, error)
}

type promptService struct {
	promptRepo repository.PromptRepository
	logger     *zap.Logger
}

func NewPromptService(promptRepo repository.PromptRepository, logger *zap.Logger) PromptService {
	return &promptService{
		promptRepo: promptRepo,
		logger:     logger,
	}
}

func (s promptService) CreateNewDescription(ctx *gin.Context, data, addons string) (string, error) {
	promptText := fmt.Sprintf("این %s دیتای یک آگهی املاک است و از نظر آگهی گذار به این دلایل %s خوب است. لطفاً ۳ تا ۵ دلیل کوتاه و جذاب استخراج کن که نشان بده چرا این ملک خوب است. \nخروجی را به صورت خیلی کوتاه بده و ساده و تبلیغاتی بنویس.", data, addons)
	res, err := s.promptRepo.Get(ctx, promptText)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s promptService) CreateAgahiNewDescription(ctx *gin.Context, input string) (string, error) {
	promptText := fmt.Sprintf(AgahiPrompt, input)
	res, err := s.promptRepo.Get(ctx, promptText)
	if err != nil {
		return "", err
	}

	startIdx := strings.Index(res, "{")
	endIdx := strings.LastIndex(res, "}")
	if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
		jsonSubstring := res[startIdx : endIdx+1]
		var optimized descRes
		if err = json.Unmarshal([]byte(jsonSubstring), &optimized); err == nil {
			formattedOutput := strings.ReplaceAll(optimized.NewDescription, "\\n", "\n")
			return formattedOutput, nil
		}
	}

	if len(res) < 1000 {
		return res, nil
	}

	promptText = fmt.Sprintf("این متن آگهی  بیشتر از ۱۰۰۰ کارکتر است. کمترین تعداد کارکتر ممکن رو حذف کن تا متن کمتر از ۱۰۰۰ کارکتر بشه. از حذف اطلاعات مهم آگهی خوداری کن. %s", res)
	res, err = s.promptRepo.Get(ctx, promptText)
	if err != nil {
		return "", err
	}

	return res, nil
}

type descRes struct {
	NewDescription string `json:"new_description"`
}
