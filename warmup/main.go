package main

import (
	"bytes"
	"chera_khube/internal/helper"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	//s := "⭐ دندانپزشکی مدرن میرداماد ⭐\n\n✅️ خدمات تخصصی و مدرن دندانپزشکی با گارانتی\n\n⚜️ اگر به دنبال تجربه‌ای متفاوت و حرفه‌ای در دندانپزشکی هستید، مرکز تخصصی ما انتخابی بی‌نظیر است. خدمات ما شامل:\n\n⏺ ارتودنسی ثابت و نامرئی با طول درمان کوتاه (تا شش ماه) و امکان پرداخت اقساط\n⏺ اصلاح و طراحی خط لبخند با کامپوزیت ونیر و لمینت سرامیکی بدون تراش\n⏺ ایمپلنت فوری و بدون جراحی\n⏺ حفظ و بازسازی ریشه‌های باقیمانده بدون نیاز به کشیدن دندان\n⏺ جراحی‌های پیچیده دندان‌های نهفته\n⏺ دندانپزشکی اطفال با روش‌های ویژه و بدون درد\n⏺ ساخت پروتزهای ثابت و متحرک بر پایه ایمپلنت یا دندان‌ها و ریشه‌های باقیمانده\n\n⭕ تجهیزات پیشرفته و مدرن:\n✅️ دستگاه بی‌حسی دیجیتال بدون درد\n✅️ قالب‌گیری دیجیتال\n✅️ استریلیزاسیون ویژه\n\n❗ ویزیت رایگان\n\n↩ قیمت خدمات: توافقی\n\n☎️ برای اطلاعات بیشتر و رزرو وقت، با ما تماس بگیرید."
	//n := 1000
	//
	//if len(s) > n {
	//	fmt.Printf("The string is longer than %d characters\n", n)
	//} else {
	//	fmt.Printf("The string is NOT longer than %d characters\n", n)
	//}
	s := "⭐ شرکت تابلوسازی و چاپ چهره ⭐\n\n✅️ عضو رسمی اتحادیه تابلو سازان با 35 سال سابقه درخشان\n✅️ شماره ثبت: 91193\n\n⭕ خدمات تخصصی و بدون واسطه:\n⏺ چاپ بنر، فلکس، مش، استیکر، سولیت\n⏺ طراحی و اجرای انواع تابلوهای برجسته:\n   ↪ حروف برجسته چنلیوم، استیل، فلزی، وکیوم، لاسوگاسی، پلاستیک، لبه سوئدی، سفارشی، پانچی\n⏺ طراحی و اجرای نما کامپوزیت\n⏺ ساخت لایت باکس و تابلو LED\n⏺ طراحی و ساخت انواع تندیس و لوح تقدیر\n⏺ برش لیزری با بالاترین کیفیت روی متریال‌های مختلف:\n   ↪ چوب، چرم، پلکسی، MDF، پارچه، فوم، نمد و ...\n\n⚜️ شعار ما: \"قیمت، صداقت، کیفیت\" – یک واقعیت است، نه فقط یک شعار!\n\n❗ ساعات کاری: از ساعت 9 صبح (تعطیلی فقط جمعه‌ها)\n\n☎️ برای اطلاعات بیشتر با ما تماس بگیرید یا به وب‌سایت ما مراجعه کنید"
	//fmt.Println(s)
	fmt.Println(len(s))
	//p := fmt.Sprintf("این متن آگهی  بیشتر از ۱۰۰۰ کارکتر است. کمترین تعداد کارکتر ممکن رو حذف کن تا متن کمتر از ۱۰۰۰ کارکتر بشه. از حذف اطلاعات مهم آگهی خوداری کن. %s", s)
	//res, err := Get(nil, p)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//fmt.Println(res)
}

const (
	BaseURL = "https://api.openai.com/v1/responses"
)

type prompt struct {
	logger *zap.Logger
	config *helper.ServiceConfig
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type OptimizedOutput struct {
	NewDescription string `json:"new_description"`
}

func Get(ctx *gin.Context, input string) (string, error) {
	//promptText := fmt.Sprintf(promptTemplate, input)

	reqPayload := ChatRequest{
		Model:       "gpt-4o",
		Messages:    []Message{{Role: "user", Content: input}},
		Temperature: 0.4,
		MaxTokens:   1000,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		//logger.Error("failed to marshal request payload", zap.Error(err))
		return "", err
	}

	//url := fmt.Sprintf("%s/chat/completions", "https://api.metisai.ir/openai/v1")
	url := fmt.Sprintf("%s/chat/completions", "https://api.avalapis.ir/v1")
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+"aa-JDpuqsYctmaNiypjpigyO1x45xn1EgfXB8hHO73T0QyGg5Zw")
	//httpReq.Header.Set("Authorization", "Bearer "+"tpsg-5k5vJuD1uWz0svsf46lpnWVqB1u0tTY")
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		//p.logger.Error("failed to send HTTP request", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API request failed: %s", string(bodyBytes))
		return "", fmt.Errorf(errMsg)
	}

	var chatResp ChatResponse
	err = json.Unmarshal(bodyBytes, &chatResp)
	if err != nil {
		return "", err
	}

	var rawOutput string
	if len(chatResp.Choices) > 0 {
		rawOutput = strings.TrimSpace(chatResp.Choices[0].Message.Content)
	}

	var optimized OptimizedOutput
	if err = json.Unmarshal([]byte(rawOutput), &optimized); err == nil {
		formattedOutput := strings.ReplaceAll(optimized.NewDescription, "\\n", "\n")
		return formattedOutput, nil
	}

	return rawOutput, nil
}

// Encode is responsible for encoding a given ID.
func Encode(id uint, alphabet, salt string, minLength int) (string, error) {

	hashData := hashids.HashIDData{
		Alphabet:  alphabet,
		Salt:      salt,
		MinLength: minLength,
	}

	h, err := hashids.NewWithData(&hashData)
	if err != nil {
		return "", errors.Wrap(err, "error on creating the hash maker")
	}

	encodedID, err := h.Encode([]int{int(id)})
	if err != nil {
		return "", errors.Wrap(err, "error on encoding hash")
	}

	return encodedID, nil
}

type res struct {
	NewDescription string `json:"new_description"`
}
