package dto

type OAuthToken struct {
	BaseUrl      string // https://api.divar.ir/oauth2/token
	Code         string
	ClientID     string
	ClientSecret string
	GrantType    string // authorization_code
	RedirectUri  string
}

type OAuthLogin struct {
	BaseUrl      string
	ResponseType string
	RedirectUri  string
	Scope        string
	State        string
	ClientId     string
}
