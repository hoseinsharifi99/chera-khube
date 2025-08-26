package api

import (
	"bytes"
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

type widgetApi struct {
	config *helper.ServiceConfig
	logger *zap.Logger
}

func NewWidgetApi(
	config *helper.ServiceConfig,
	logger *zap.Logger,
) repository.WidgetRepository {
	return &widgetApi{
		config: config,
		logger: logger,
	}
}

func (i widgetApi) Delete(postToken, accessToken string) error {
	endPoint := deleteURL(i.config.Divar.DeleteWidget, postToken)

	method := "DELETE"
	req, err := http.NewRequest(method, endPoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", i.config.CarDivar.ApiKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response res
		err = json.NewDecoder(resp.Body).Decode(&response)
		return errors.New(response.Message)
	}

	return nil
}

func (i widgetApi) DeleteByPostToken(postToken, accessToken string) error {
	endPoint := deletePostURL(i.config.CarDivar.DeleteWidget, postToken)

	method := "DELETE"
	req, err := http.NewRequest(method, endPoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", i.config.CarDivar.ApiKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response res
		err = json.NewDecoder(resp.Body).Decode(&response)
		return errors.New(response.Message)
	}

	return nil
}

func (i widgetApi) Send(widget *model.DivarWidget, postToken, accessToken string) error {
	endPoint := addOnsURL(i.config.Divar.AddOns, postToken)

	body, err := json.Marshal(widget)
	if err != nil {
		return err
	}

	method := "POST"
	req, err := http.NewRequest(method, endPoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("x-api-key", i.config.CarDivar.ApiKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response res
		err = json.NewDecoder(resp.Body).Decode(&response)
		return errors.New(response.Message)
	}

	log.Println("DIVAR ADD WIDGET RESPONSE")

	return nil
}

func addOnsURL(endpoint string, postToken string) string {
	return strings.Replace(endpoint, "{{post_token}}", postToken, 1)
}

func deleteURL(endpoint string, id string) string {
	return strings.Replace(endpoint, "{{post_token}}", id, 1)
}

func deletePostURL(endpoint string, id string) string {
	return strings.Replace(endpoint, "{token}", id, 1)
}

type res struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type okRes struct {
	Id string `json:"id"`
}
