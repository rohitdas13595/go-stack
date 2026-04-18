package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the default JWT payload.
type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

// Issue creates a signed JWT for userID.
func Issue(secret []byte, userID string, ttl time.Duration) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("auth: empty secret")
	}
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

// Parse validates token and returns claims.
func Parse(secret []byte, tokenString string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := t.Claims.(*Claims); ok && t.Valid {
		return c, nil
	}
	return nil, fmt.Errorf("auth: invalid token")
}
