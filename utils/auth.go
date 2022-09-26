package utils

import (
	"errors"
	"go-rest/settings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const oneMonth = 43800

type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(email string, id int64) string {
	expirationTime := time.Now().Add(oneMonth * 3 * time.Minute).Unix()
	claims := &Claims{
		ID:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(settings.Secret))
	return tokenString
}

func GetClaims(token string) (Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) { return []byte(settings.Secret), nil })
	if err != nil {
		return *claims, err
	}
	return *claims, nil
}

func VerifyToken(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) { return []byte(settings.Secret), nil })
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", errors.New("token is invalid")
	}
	expiry := time.Unix(claims.ExpiresAt, 0)
	if time.Until(expiry) < 1*time.Second {
		return "", errors.New("token's time expired")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := _token.SignedString([]byte(settings.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
