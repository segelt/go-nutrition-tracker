package internal

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtManager struct {
	jwt_secret []byte
}

func New(JWT_SECRET string) *JwtManager {
	return &JwtManager{
		jwt_secret: []byte(JWT_SECRET),
	}
}

type JWTClaim struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func (m JwtManager) GenerateToken(userId string, username string) (generatedToken string, err error) {
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
	generatedToken, err = token.SignedString(m.jwt_secret)
	return generatedToken, err
}

func (m JwtManager) ValidateToken(encodedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			// return []byte(jwtKey), nil
			return m.jwt_secret, nil
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
