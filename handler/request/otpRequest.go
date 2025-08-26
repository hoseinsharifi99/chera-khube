package request

type OtpRequest struct {
	Barcode         *string `json:"barcode"`
	DocumentNumber  *string `json:"document_number"`
	BuyerNationalID string  `json:"buyer_national_id"`
}

type InquiryRequest struct {
	BuyerNationalID string  `json:"buyer_national_id"`
	Barcode         *string `json:"barcode"`
	DocumentNumber  *string `json:"document_number"`
	OtpCode         *string `json:"otp_code"`
	OtpRef          *string `json:"otp_reference_number"`
	Trace           *string `json:"trace_number"`
}
