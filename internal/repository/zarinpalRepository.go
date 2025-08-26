package repository

import (
	"chera_khube/internal/dto"
	"chera_khube/internal/response"
)

type ZarinpalRepository interface {
	NewPaymentRequest(d dto.PaymentRequestDto) (*dto.PaymentResponseDto, error)
	PaymentVerification(d dto.PaymentVerificationDto) (*dto.PaymentVerificationResponseDto, error)
	UnverifiedTransactions() (authorities []response.UnverifiedAuthority, statusCode int, err error)
	RefreshAuthority(authority string, expire int) (statusCode int, err error)
}
