package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go-gin-admin/model"
	"time"
)

var jwtKey = []byte(viper.GetString("Jwt.key"))

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func ReleaseToken(user model.AdminUsers) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "D A I",
			Subject:   "User token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
