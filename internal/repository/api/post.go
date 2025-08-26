package api

import (
	"chera_khube/internal/helper"
	"chera_khube/internal/model"
	"chera_khube/internal/repository"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
)

type postApi struct {
	config *helper.ServiceConfig
	logger *zap.Logger
}

func NewPostApi(
	config *helper.ServiceConfig,
	logger *zap.Logger,
) repository.PostApiRepo {
	return &postApi{
		config: config,
		logger: logger,
	}
}

func (i postApi) Get(token string) (*model.Post, error) {
	endPoint := getPostUrl(i.config.Divar.GetPost, token)

	method := "GET"
	req, err := http.NewRequest(method, endPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", i.config.Divar.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("method:getPostDivar err status_code:", resp.StatusCode)
		return nil, errors.New(resp.Status)
	}

	var response getPostResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.New("error while decode response")
	}
	return response.toPostModel()
}

func (i postApi) GetByAll(token string) (*model.Post, error) {
	endPoint := getPostUrl(i.config.Divar.GetPost, token)

	method := "GET"
	req, err := http.NewRequest(method, endPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", i.config.AgahiPlus.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("method:getPostDivar err status_code:", resp.StatusCode)
		return nil, errors.New(resp.Status)
	}

	var response getPostResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.New("error while decode response")
	}
	return response.toPostModel()
}

func getPostUrl(endpoint string, postToken string) string {
	return strings.Replace(endpoint, "{{token}}", postToken, 1)
}

type getPostResponse struct {
	Token       string                 `json:"token"`
	Category    string                 `json:"category"`
	City        string                 `json:"city"`
	District    string                 `json:"district"`
	ChatEnabled bool                   `json:"chat_enabled"`
	Data        map[string]interface{} `json:"data"` // Dynamic data field
}

func (r getPostResponse) toPostModel() (*model.Post, error) {
	// Marshal entire response into JSON for Data field
	data, err := json.Marshal(r.Data)
	if err != nil {
		return nil, errors.New("error while marshalling response data")
	}

	// Extract title and description dynamically
	title := ""
	if val, ok := r.Data["title"].(string); ok {
		title = val
	}

	// Populate Post model
	return &model.Post{
		Token:       r.Token,
		Data:        string(data),
		IsConnected: false,
		Title:       title,
		Category:    r.Category,
		City:        r.City,
		District:    r.District,
	}, nil
}

//build bgir koskesh
