package session

import (
	"time"

	"github.com/dattito/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store

func Init() {
	SessionStore = session.New(session.Config{
		CookieHTTPOnly: config.AUTHORIZATION_COOKIE_HTTPONLY,
		CookieSecure:   config.AUTHORIZATION_COOKIE_SECURE,
		CookieSameSite: config.AUTHORIZATION_COOKIE_SAMESITE,
		CookieDomain:   config.AUTHORIZATION_COOKIE_DOMAIN,
		Expiration:     10 * time.Minute,
	})
}
