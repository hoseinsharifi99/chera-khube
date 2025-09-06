package api

import (
	"bytes"
	"chera_khube/internal/dto"
	"chera_khube/internal/repository"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
)

type oAuthRepository struct {
	logger *zap.Logger
}

func NewOAuthRepository(logger *zap.Logger) repository.OAuthRepository {
	return &oAuthRepository{logger: logger}
}

func (o *oAuthRepository) GetToken(d dto.OAuthToken) (*dto.AccessTokenResponse, error) {
	apiUrl := d.BaseUrl
	method := "POST"

	data := url.Values{}
	data.Set("code", d.Code)
	data.Set("client_id", d.ClientID)
	data.Set("client_secret", d.ClientSecret)
	data.Set("grant_type", d.GrantType)
	data.Set("redirect_uri", d.RedirectUri)

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, bytes.NewBufferString(data.Encode()))
	log.Println("url:", d.BaseUrl)
	log.Println("code:", d.Code)
	log.Println("client_id:", d.ClientID)
	log.Println("client_secret:", d.ClientSecret)
	log.Println("grant_type:", d.GrantType)
	log.Println("redirect_uri:", d.RedirectUri)

	if err != nil {
		o.logger.Warn("Failed to create request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		o.logger.Warn("Failed to get token", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		o.logger.Warn("Error reading body", zap.Error(err))
		return nil, err
	}

	o.logger.Info("Successfully get token", zap.String("token", string(body)))

	var response dto.AccessTokenResponse
	err = json.Unmarshal(body, &response)

	return &response, nil
}

func (o *oAuthRepository) GetPhoneNumber(d dto.PhoneNumber) (*dto.PhoneNumberResponse, error) {
	url := d.BaseUrl
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		o.logger.Warn("Failed to create request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+d.AccessToken)
	req.Header.Add("x-api-key", d.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		o.logger.Warn("Failed to get token", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		o.logger.Warn("Error reading body", zap.Error(err))
		return nil, err
	}

	o.logger.Info("Successfully get phone number", zap.String("phone number", string(body)))

	var response dto.PhoneNumberResponse
	err = json.Unmarshal(body, &response)

	return &response, nil
}
