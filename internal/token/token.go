package token

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	*jwt.RegisteredClaims
}

func GenerateToken(secret, userID string) (string, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		UserID: userID,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:   "golang-be-template",
			IssuedAt: &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
		},
	})

	return token.SignedString(decodedSecret)
}

func ValidateToken(secret, tokenString string) (*TokenClaims, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return decodedSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*TokenClaims), nil
}
