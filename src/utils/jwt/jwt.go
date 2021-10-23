package jwt

import (
	"fmt"
	"time"

	"github.com/dattito/purrmannplus-backend/config"
	"github.com/golang-jwt/jwt/v4"
)

// Creates a new JWT token for the given user including the account_id
func NewAccountIdToken(accountId string) (string, time.Time, error) {
	expires := time.Now().Add(time.Duration(config.AUTHORIZATION_EXPIRATION_TIME) * time.Second)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["account_id"] = accountId
	claims["exp"] = expires.Unix()

	t, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", time.Time{}, err
	}

	return t, expires, nil
}

// Creates a new short living JWT token for the user including the account_id and the phone_number
func NewAccountIdPhoneNumberToken(accountId, phone_number string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["account_id"] = accountId
	claims["phone_number"] = phone_number
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	t, err := token.SignedString([]byte(config.JWT_SHORTLIVING_SECRET))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Returns the account_id from the given token
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

// Returns the account_id and phone_number from the given token
func ParseAccountIdPhoneNumberToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWT_SHORTLIVING_SECRET), nil
	})
	if err != nil {
		return "", "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["account_id"].(string), claims["phone_number"].(string), nil
}
