package api

import (
	"bytes"
	"chera_khube/handler/request"
	"chera_khube/internal/constant"
	"chera_khube/internal/dto"
	"chera_khube/internal/repository"
	"chera_khube/internal/response"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Zarinpal is the base struct for zarinpal payment
// gateway, one shall not create or manipulate instances
// if this struct manually and just use provided methods
// to week with it.
type zarinpal struct {
	MerchantID      string
	Sandbox         bool
	APIEndpoint     string
	PaymentEndpoint string
}

func NewZarinpal(merchantID string, sandbox bool) repository.ZarinpalRepository {
	if len(merchantID) != 36 {
		return nil
	}
	apiEndPoint := constant.ZarinpalApiEndPoint
	paymentEndpoint := constant.ZarinpalPaymentEndpoint
	if sandbox == true {
		apiEndPoint = constant.ZarinpalSandboxApiEndPoint
		paymentEndpoint = constant.ZarinpalSandboxPaymentEndpoint
	}
	return &zarinpal{
		MerchantID:      merchantID,
		Sandbox:         sandbox,
		APIEndpoint:     apiEndPoint,
		PaymentEndpoint: paymentEndpoint,
	}
}

func (z *zarinpal) NewPaymentRequest(d dto.PaymentRequestDto) (*dto.PaymentResponseDto, error) {
	meta := request.MetaData{
		Email:  d.Email,
		Mobile: d.PhoneNumber,
	}

	paymentRequest := request.PaymentRequestReqBody{
		MerchantID:  z.MerchantID,
		Amount:      d.Amount,
		CallbackURL: d.CallbackUrl,
		Description: d.Description,
		MetaData:    meta,
	}
	var resp response.PaymentRequestResp
	err := z.request("request.json", &paymentRequest, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.Code == 100 {
		return &dto.PaymentResponseDto{
			PaymentUrl: z.PaymentEndpoint + resp.Data.Authority,
			Authority:  resp.Data.Authority,
			StatusCode: resp.Data.Code,
		}, nil
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Code))
	}
	return nil, err
}

func (z *zarinpal) PaymentVerification(d dto.PaymentVerificationDto) (*dto.PaymentVerificationResponseDto, error) {
	paymentVerification := request.PaymentVerificationReqBody{
		MerchantID: z.MerchantID,
		Amount:     d.Amount,
		Authority:  d.Authority,
	}

	var resp response.PaymentVerificationResp
	err := z.request("verify.json", &paymentVerification, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.Code == 100 {
		return &dto.PaymentVerificationResponseDto{
			Verified:   true,
			RefID:      strconv.Itoa(resp.Data.RefId),
			StatusCode: resp.Data.Code,
		}, nil
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Code))
	}

	return nil, err
}

func (z *zarinpal) UnverifiedTransactions() (authorities []response.UnverifiedAuthority, statusCode int, err error) {
	unverifiedTransactions := request.UnverifiedTransactionsReqBody{
		MerchantID: z.MerchantID,
	}

	var resp response.UnverifiedTransactionsResp
	err = z.request("UnverifiedTransactions.json", &unverifiedTransactions, &resp)
	if err != nil {
		return
	}

	if resp.Status == 100 {
		statusCode = resp.Status
		authorities = resp.Authorities
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

func (z *zarinpal) RefreshAuthority(authority string, expire int) (statusCode int, err error) {
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}
	if expire < 1800 {
		err = errors.New("expire must be at least 1800")
		return
	} else if expire > 3888000 {
		err = errors.New("expire must not be greater than 3888000")
		return
	}

	refreshAuthority := request.RefreshAuthorityReqBody{
		MerchantID: z.MerchantID,
		Authority:  authority,
		ExpireIn:   expire,
	}
	var resp response.RefreshAuthorityResp
	err = z.request("RefreshAuthority.json", &refreshAuthority, &resp)
	if err != nil {
		return
	}
	if resp.Status == 100 {
		statusCode = resp.Status
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

func (z *zarinpal) request(method string, data interface{}, res interface{}) error {
	reqBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Println(z.APIEndpoint + method)
	req, err := http.NewRequest("POST", z.APIEndpoint+method, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		err = errors.New("zarinpal invalid json response")
		return err
	}
	return nil
}
