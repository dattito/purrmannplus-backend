package jwt

import (
	"time"

	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/golang-jwt/jwt/v4"
)

func NewAccountLoginToken(accountId string) (string, error) {
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
