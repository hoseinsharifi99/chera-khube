package api

import (
	"bytes"
	"chera_khube/internal/helper"
	"chera_khube/internal/repository"
	"chera_khube/internal/response"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func (i divarApi) EditDescription(token, accessToken, description string) error {
	endPoint := getEditUrl(i.config.Divar.EditPost, token)
	log.Println("endpoint: ", endPoint)
	requestModel := &reqModel{
		Description: description,
	}
	body, err := json.Marshal(requestModel)
	if err != nil {
		return err
	}

	method := "PUT"
	req, err := http.NewRequest(method, endPoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", i.config.Divar.ApiKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("method:edit DivarPost err status_code:", resp.StatusCode)
		log.Println("error on divar edit:", string(b), "ENDPOINT:", endPoint, "BODY:", string(body), "ACCESS_TOKEN:", accessToken)
		return errors.New(resp.Status)
	}

	fmt.Println(string(b))
	return nil
}

func (i divarApi) EditAllDescription(token, accessToken, description string) error {
	endPoint := getEditUrl(i.config.AgahiPlus.EditPost, token)
	log.Println("endpoint: ", endPoint)
	requestModel := &reqModel{
		Description: description,
	}
	body, err := json.Marshal(requestModel)
	if err != nil {
		return err
	}

	method := "PUT"
	req, err := http.NewRequest(method, endPoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", i.config.AgahiPlus.ApiKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("method:edit DivarPost err status_code:", resp.StatusCode)
		log.Println("error on divar edit:", string(b), "ENDPOINT:", endPoint, "BODY:", string(body), "ACCESS_TOKEN:", accessToken)
		return errors.New(resp.Status)
	}

	fmt.Println(string(b))
	return nil
}

func getEditUrl(endpoint string, postToken string) string {
	return strings.Replace(endpoint, "{{token}}", postToken, 1)
}

type reqModel struct {
	Description string `json:"description"`
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
