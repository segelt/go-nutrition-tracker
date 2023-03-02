package internal

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("ToBeConfiguredFromEnv")

type JWTClaim struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(userId string, username string) (generatedToken string, err error) {
	expirationTime := time.Now().Add(time.Hour * 48)
	claims := &JWTClaim{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        userId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	generatedToken, err = token.SignedString(jwtKey)
	return generatedToken, err
}

func ValidateToken(encodedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, err

}
