package api

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/repository"
	"chera_khube/internal/response"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type divarApi struct {
	config *helper.ServiceConfig
	logger *zap.Logger
}

func NewDivarApi(
	config *helper.ServiceConfig,
	logger *zap.Logger,
) repository.DivarRepository {
	return &divarApi{
		config: config,
		logger: logger,
	}
}

func (i divarApi) GetPostTokens(url, apikey, accessToken string) (*response.GetPostsResponse, error) {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", apikey)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	var resp response.GetPostsResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &resp, nil
}
