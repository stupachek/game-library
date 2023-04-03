package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrParseClaim error = errors.New("can't parse claim")
	ErrOutdated   error = errors.New("token is outdated")
)

var key = []byte("supersecretkey")

type JWTClaim struct {
	Email string `json:"id"`
	jwt.StandardClaims
}

func NewJWT(email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func ValidateToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) { return []byte(key), nil })
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", ErrParseClaim
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return "", ErrOutdated
	}
	return claims.Email, nil
}
