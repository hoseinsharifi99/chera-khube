package response

type PaymentRequestResp struct {
	Data struct {
		Authority string `json:"authority"`
		Fee       int    `json:"fee"`
		FeeType   string `json:"fee_type"`
		Code      int    `json:"code"`
		Message   string `json:"message"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

type PaymentVerificationResp struct {
	Data struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CardHash string `json:"card_hash"`
		CardPan  string `json:"card_pan"`
		RefId    int    `json:"ref_id"`
		FeeType  string `json:"fee_type"`
		Fee      int    `json:"fee"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

// UnverifiedAuthority is the base struct for Authorities in unverifiedTransactionsResp
type UnverifiedAuthority struct {
	Authority   string
	Amount      int
	Channel     string
	CallbackURL string
	Referer     string
	Email       string
	CellPhone   string
	Date        string // ToDo Check type to be date
}

type UnverifiedTransactionsResp struct {
	Status      int
	Authorities []UnverifiedAuthority
}

type RefreshAuthorityResp struct {
	Status int
}
