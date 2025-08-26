package repository

import "chera_khube/internal/response"

type DivarRepository interface {
	EditDescription(token, accessToken, description string) error
	EditAllDescription(token, accessToken, description string) error
	GetPostTokens(url, apikey, accessToken string) (*response.GetPostsResponse, error)
}
