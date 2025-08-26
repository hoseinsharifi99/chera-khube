package middlewares

import (
	"chera_khube/internal/constant/errorKey"
	"chera_khube/internal/constant/messageKey"
	"chera_khube/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	PHONE_NUMBER = "phoneNumber"
	PASSWORD     = "password"
)

func extractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func Jwt(e *helper.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" || tokenString == "undefined" {
			c.JSON(http.StatusUnauthorized, gin.H{messageKey.ERROR: errorKey.UNAUTHORIZED})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &helper.SignedDetail{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(e.JWT.Secret), nil
		})

		if claims, ok := token.Claims.(*helper.SignedDetail); ok && token.Valid {
			if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
				c.JSON(http.StatusForbidden, gin.H{messageKey.ERROR: errorKey.INVALID_TOKEN})
				c.Abort()
				return
			}
			c.Set(PHONE_NUMBER, claims.PhoneNumber)
			c.Set(PASSWORD, claims.Password)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{messageKey.ERROR: err.Error()})
			c.Abort()
			return
		}
	}
}

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, authOK := c.Request.BasicAuth()
		if !authOK {
			c.JSON(http.StatusUnauthorized, gin.H{messageKey.ERROR: errorKey.UNAUTHORIZED})
			c.Abort()
			return
		}
		hasWhitespace := regexp.MustCompile(`\s`).MatchString(password)
		if hasWhitespace || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{messageKey.ERROR: errorKey.PASSWORD_IS_NULL})
			c.Abort()
			return
		}
		c.Set(PHONE_NUMBER, username)
		c.Set(PASSWORD, password)
	}
}
