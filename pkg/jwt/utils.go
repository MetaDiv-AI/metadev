package jwt

import (
	"fmt"
	"maps"

	"github.com/golang-jwt/jwt"
)

func jwtTokenToClaims(token *jwt.Token) (Claims, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("jwt error - invalid claims type")
	}
	claims := ClaimsImpl{
		values: make(map[string]any),
	}
	maps.Copy(claims.values, mapClaims)
	return &claims, nil
}
