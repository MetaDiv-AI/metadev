package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"maps"

	"github.com/golang-jwt/jwt"
)

// NewRsaEncoder creates a new RsaEncoder that can be used to encode and decode JWT tokens
// using an RSA private or public key
func NewRsaEncoder(privateOrPublicPem string) RsaEncoder {
	return &RsaEncoderImpl{
		pem: privateOrPublicPem,
	}
}

// RsaEncoder is the interface for encoding and decoding JWT tokens using an RSA private or public key
type RsaEncoder interface {
	// Encode encodes the claims into a JWT token
	Encode(claims Claims) (string, error)
	// Decode decodes the JWT token into claims
	Decode(token string) (Claims, error)
	// PublicDecode decodes the JWT token into claims using the public key
	PublicDecode(token string) (Claims, error)
}

// RsaEncoderImpl is the implementation of the RsaEncoder interface
type RsaEncoderImpl struct {
	pem string
}

// Encode encodes the claims into a JWT token
func (r *RsaEncoderImpl) Encode(claims Claims) (string, error) {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(r.pem))
	if err != nil {
		return "", fmt.Errorf("jwt error - parsing private key PEM: %w", err)
	}
	jwtClaims := jwt.MapClaims{}
	maps.Copy(jwtClaims, claims.Values())
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims).SignedString(parsedKey)
	if err != nil {
		return "", fmt.Errorf("jwt error - signing token: %w", err)
	}
	return tokenString, nil
}

// Decode decodes the JWT token into claims
func (r *RsaEncoderImpl) Decode(token string) (Claims, error) {
	parsedKey, err := r.parsePublicKey()
	if err != nil {
		return nil, err
	}
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return parsedKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt error - parsing token: %w", err)
	}
	return jwtTokenToClaims(jwtToken)
}

// PublicDecode decodes the JWT token into claims using the public key
func (r *RsaEncoderImpl) PublicDecode(token string) (Claims, error) {
	parsedKey, err := r.parsePublicKey()
	if err != nil {
		return nil, err
	}
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return parsedKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt error - parsing token: %w", err)
	}
	return jwtTokenToClaims(jwtToken)
}

// parsePublicKey attempts to parse the public key in multiple formats
func (r *RsaEncoderImpl) parsePublicKey() (interface{}, error) {
	// Try PKIX format first (newer format)
	if parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(r.pem)); err == nil {
		return parsedKey, nil
	}
	
	// Try PKCS1 format (older format)
	if parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(r.pem)); err == nil {
		return parsedKey, nil
	}
	
	// Try parsing as generic public key
	block, _ := pem.Decode([]byte(r.pem))
	if block == nil {
		return nil, fmt.Errorf("jwt error - failed to decode PEM block")
	}
	
	switch block.Type {
	case "PUBLIC KEY":
		parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("jwt error - parsing PKIX public key: %w", err)
		}
		return parsedKey, nil
	case "RSA PUBLIC KEY":
		parsedKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("jwt error - parsing PKCS1 public key: %w", err)
		}
		return parsedKey, nil
	default:
		return nil, fmt.Errorf("jwt error - unsupported public key type: %s", block.Type)
	}
}
