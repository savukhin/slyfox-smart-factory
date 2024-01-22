package service

import (
	"errors"
	"eventsproxy/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	User domain.User
	jwt.RegisteredClaims
}

func (svc *proxyService) generateJWT(user domain.User) (string, error) {
	claims := &jwtClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(svc.jwtDuration)),
			Issuer:    user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(svc.jwtSecret)
	return ss, err
}

func (svc *proxyService) parseJWT(tokenString string) (domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return svc.jwtSecret, nil
	})

	if err != nil {
		return domain.User{}, err
	} else if claims, ok := token.Claims.(*jwtClaims); ok {
		return claims.User, nil
	} else {
		fmt.Println(token.Claims)
		return domain.User{}, errors.New("unknown claims type")
	}
}
