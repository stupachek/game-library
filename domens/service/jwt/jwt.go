package jwt

import (
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrParseClaim error = errors.New("can't parse claim")
	ErrOutdated   error = errors.New("token is outdated")
	Public        ed25519.PublicKey
	Private       ed25519.PrivateKey
)

type JWTClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func NewJWT(ID string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		ID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(Private)
}

func ValidateToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) { return Public, nil })
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
	return claims.ID, nil
}
