package dto

type PhoneNumberResponse struct {
	PhoneNumbers []string `json:"phone_numbers"`
	PhoneNumber  string   `json:"phone_number"`
}
