package jwt

import (
	"fmt"
	"maps"

	"github.com/golang-jwt/jwt"
)

// NewSecretEncoder creates a new SecretEncoder that can be used to encode and decode JWT tokens
// using a secret key
func NewSecretEncoder(secret string) SecretEncoder {
	return &SecretEncoderImpl{
		secret: secret,
	}
}

// SecretEncoder is the interface for encoding and decoding JWT tokens using a secret key
type SecretEncoder interface {
	// Encode encodes the claims into a JWT token
	Encode(claims Claims) (string, error)
	// Decode decodes the JWT token into claims
	Decode(token string) (Claims, error)
}

// SecretEncoderImpl is the implementation of the SecretEncoder interface
type SecretEncoderImpl struct {
	secret string
}

// Encode encodes the claims into a JWT token
func (s *SecretEncoderImpl) Encode(claims Claims) (string, error) {
	jwtClaims := jwt.MapClaims{}
	maps.Copy(jwtClaims, claims.Values())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("jwt error - signing token: %w", err)
	}
	return tokenString, nil
}

// Decode decodes the JWT token into claims
func (s *SecretEncoderImpl) Decode(token string) (Claims, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt error - parsing token: %w", err)
	}
	return jwtTokenToClaims(jwtToken)
}
