package token

import (
	"net/http"
	"time"
	"user-service/config"
	"user-service/db/user"
)

func ToCookie(name string, token string, expires time.Duration) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  time.Now().Add(expires),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   config.GO_ENV == "production",
	}
}

func CreateAccessCookie(user user.User, rememberMe bool) http.Cookie {
	e := AccessExpiry
	c := CreateAccess(user)
	if !rememberMe {
		e = 0
	}

	return ToCookie(config.ACCESS_TOKEN_KEY, c, e)
}

func CreateRefreshCookie(user user.User, rememberMe bool) http.Cookie {
	e := RefreshExpiry
	c := CreateRefresh(user)
	if !rememberMe {
		e = 0
	}

	return ToCookie(config.REFRESH_TOKEN_KEY, c, e)
}
