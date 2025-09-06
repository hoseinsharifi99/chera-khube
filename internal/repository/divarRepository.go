package repository

import "chera_khube/internal/response"

type DivarRepository interface {
	GetPostTokens(url, apikey, accessToken string) (*response.GetPostsResponse, error)
}
