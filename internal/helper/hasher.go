package helper

import (
	"crypto/sha256"
	"fmt"
)

func GetPassword(password string, secretKey string) string {
	sha := sha256.Sum256([]byte(password + secretKey))
	return fmt.Sprintf("%x", sha)
}
