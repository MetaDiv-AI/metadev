package jwt

import (
	"fmt"
	"maps"

	"github.com/golang-jwt/jwt"
)

// NewUnverifiedEncoder creates a new UnverifiedEncoder that can be used to encode and decode JWT tokens
// using the none signing method
func NewUnverifiedEncoder() UnverifiedEncoder {
	return &UnverifiedEncoderImpl{}
}

// UnverifiedEncoder is the interface for encoding and decoding JWT tokens using the none signing method
type UnverifiedEncoder interface {
	// Encode encodes the claims into a JWT token
	Encode(claims Claims) (string, error)
	// Decode decodes the JWT token into claims
	Decode(token string) (Claims, error)
}

// UnverifiedEncoderImpl is the implementation of the UnverifiedEncoder interface
type UnverifiedEncoderImpl struct{}

// Encode encodes the claims into a JWT token
func (u *UnverifiedEncoderImpl) Encode(claims Claims) (string, error) {
	jwtClaims := jwt.MapClaims{}
	maps.Copy(jwtClaims, claims.Values())
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwtClaims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		return "", fmt.Errorf("jwt error - creating unverified token: %w", err)
	}
	return tokenString, nil
}

// Decode decodes the JWT token into claims
func (u *UnverifiedEncoderImpl) Decode(token string) (Claims, error) {
	jwtToken, err := jwt.Parse(token, nil)
	if err != nil && err.Error() != "token is unverifiable: no keyfunc was provided" && err.Error() != "no Keyfunc was provided." {
		return nil, fmt.Errorf("jwt error - parsing unverified token: %w", err)
	}
	return jwtTokenToClaims(jwtToken)
}
