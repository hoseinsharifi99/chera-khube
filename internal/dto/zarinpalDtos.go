package dto

type PaymentRequestDto struct {
	Amount      int
	CallbackUrl string
	Description string
	Email       string
	PhoneNumber string
}

type PaymentResponseDto struct {
	PaymentUrl string
	Authority  string
	StatusCode int
}

type PaymentVerificationDto struct {
	Amount    int
	Authority string
}

type PaymentVerificationResponseDto struct {
	Verified    bool
	RefID       string
	StatusCode  int
	ServiceName string
}
