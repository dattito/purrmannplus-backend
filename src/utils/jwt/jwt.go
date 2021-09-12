package jwt

import (
	"fmt"
	"time"

	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/golang-jwt/jwt/v4"
)

func NewAccountIdToken(accountId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["account_id"] = accountId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	t, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return t, nil
}

func NewAccountIdPhoneNumberToken(accountId, phone_number string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["account_id"] = accountId
	claims["phone_number"] = phone_number
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	t, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ParseAccountIdToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["account_id"].(string), nil
}

func ParseAccountIdPhoneNumberToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		return "", "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["account_id"].(string), claims["phone_number"].(string), nil
}
