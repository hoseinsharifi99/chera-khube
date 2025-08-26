package request

type PaymentRequestReqBody struct {
	MerchantID  string   `json:"merchant_id"`
	Amount      int      `json:"amount"`
	CallbackURL string   `json:"callback_url"`
	Description string   `json:"description"`
	MetaData    MetaData `json:"metadata"`
}

type MetaData struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

type PaymentVerificationReqBody struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
	Amount     int    `json:"amount"`
}
type UnverifiedTransactionsReqBody struct {
	MerchantID string `json:"merchant_id"`
}

type RefreshAuthorityReqBody struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
	ExpireIn   int    `json:"expire_in"`
}
