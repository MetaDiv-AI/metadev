package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// NewRsaGenerator creates a new RsaGenerator that can be used to generate RSA key pairs
func NewRsaGenerator() RsaGenerator {
	return &RsaGeneratorImpl{}
}

// RsaGenerator is the interface for generating RSA key pairs
type RsaGenerator interface {
	// GenerateKeyPair generates an RSA key pair in PEM format
	GenerateKeyPair() (privateKey string, publicKey string, err error)
}

type RsaGeneratorImpl struct{}

// GenerateKeyPair generates an RSA key pair in PEM format
func (r *RsaGeneratorImpl) GenerateKeyPair() (privateKeyPem string, publicKeyPem string, err error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("jwt error - generating RSA key pair: %w", err)
	}

	privateKeyPem = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("jwt error - marshaling public key: %w", err)
	}

	publicKeyPem = string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	))
	return privateKeyPem, publicKeyPem, nil
}
