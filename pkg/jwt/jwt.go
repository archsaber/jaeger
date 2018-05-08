package jwt

import (
	"errors"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func CheckTokenValidity(token, secret string) (map[string]interface{}, error) {
	if len(token) < 6 || strings.ToUpper(token[0:6]) != "BEARER" {
		return nil, errors.New("Unuthorizated access, this event will be logged and reported")
	}
	token = token[7:]
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}
