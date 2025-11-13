package token

import (
	"time"
	"user-service/config"
	"user-service/db/user"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Method            = jwt.SigningMethodHS256
	AccessExpiry      = 5 * time.Minute                    // 5 min
	SoftRefreshExpiry = 7 * 24 * time.Hour                 // 1 week
	RefreshExpiry     = SoftRefreshExpiry + 3*24*time.Hour // 2 weeks and 3 days
)

type Claims struct {
	ID      string    `json:"id"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
	Expires time.Time `json:"expires"`
	jwt.RegisteredClaims
}

func CreateTokenClaim(user user.User, expires time.Duration, softExpires ...time.Duration) Claims {
	se := expires
	if len(softExpires) > 0 {
		se = softExpires[0];
	}
	return Claims{
		ID:      user.ID,
		Email:   user.Email,
		Role:    user.Role,
		Expires: time.Now().Add(se),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
		},
	}
}

func CreateAccess(user user.User) string {
	token := jwt.NewWithClaims(Method, CreateTokenClaim(user, AccessExpiry))
	tstr, _ := token.SignedString(config.ACCESS_TOKEN_SECRET)
	return tstr
}

func CreateRefresh(user user.User) string {
	token := jwt.NewWithClaims(Method, CreateTokenClaim(user, RefreshExpiry, SoftRefreshExpiry))
	tstr, _ := token.SignedString(config.REFRESH_TOKEN_SECRET)
	return tstr
}

func ParseAccess(token string) (Claims, bool) {
	claim := Claims{}
	t, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
		if t.Method != Method {
			return nil, jwt.ErrTokenUnverifiable
		}
		return config.ACCESS_TOKEN_SECRET, nil
	})
	if err != nil {
		return claim, false
	}
	return claim, t.Valid
}

func ParseRefresh(token string) (Claims, bool) {
	claim := Claims{}
	t, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
		if t.Method != Method {
			return nil, jwt.ErrTokenUnverifiable
		}
		return config.REFRESH_TOKEN_SECRET, nil
	})

	if err != nil {
		return claim, false
	}
	return claim, t.Valid
}
