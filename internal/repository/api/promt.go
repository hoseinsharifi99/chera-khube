package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"chera_khube/internal/helper"
	"chera_khube/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	BaseURL = "https://api.avalapis.ir/v1"
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

func (p prompt) Get(ctx *gin.Context, input string) (string, error) {
	//promptText := fmt.Sprintf(promptTemplate, input)

	reqPayload := ChatRequest{
		Model:       "gpt-4o",
		Messages:    []Message{{Role: "user", Content: input}},
		Temperature: 0.4,
		MaxTokens:   1000,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		p.logger.Error("failed to marshal request payload", zap.Error(err))
		return "", err
	}

	url := fmt.Sprintf("%s/chat/completions", p.config.App.LLMUrl)
	//url := fmt.Sprintf("%s/chat/completions", "https://api.metisai.ir/openai/v1")
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		p.logger.Error("failed to create HTTP request", zap.Error(err))
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.config.App.LLMApiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		p.logger.Error("failed to send HTTP request", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("failed to read response body", zap.Error(err))
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API request failed: %s", string(bodyBytes))
		p.logger.Error(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	var chatResp ChatResponse
	err = json.Unmarshal(bodyBytes, &chatResp)
	if err != nil {
		p.logger.Error("failed to unmarshal chat response", zap.Error(err))
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

func NewPromptApi(logger *zap.Logger, config *helper.ServiceConfig) repository.PromptRepository {
	return &prompt{
		logger: logger,
		config: config,
	}
}
