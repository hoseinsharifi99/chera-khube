package response

type GetPostsResponse struct {
	Posts []PostResponse `json:"posts"`
}

type PostResponse struct {
	Token    string   `json:"token"`
	Title    string   `json:"title"`
	Images   []string `json:"images"`
	Category string   `json:"category"`
}
