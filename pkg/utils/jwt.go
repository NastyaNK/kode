package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"myproject/pkg/config"
	"time"
)

type JwtUtils struct {
	hmacSecret []byte
	duration   int
}

func NewJwtUtils(config *config.Application) *JwtUtils {
	return &JwtUtils{
		hmacSecret: []byte(config.Auth.SecretKey),
		duration:   config.Auth.Duration,
	}
}

func (service *JwtUtils) GenerateToken(data map[string]interface{}) (string, error) {
	claims := jwt.MapClaims(data)
	claims["exp"] = time.Now().AddDate(0, 0, service.duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.hmacSecret)
}

func (service *JwtUtils) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return service.hmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token structure")
	}
	return claims, err
}
