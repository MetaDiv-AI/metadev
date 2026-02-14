package types

import (
	"os"
	"path/filepath"
	"time"

	jwtPkg "github.com/MetaDiv-AI/metadev/pkg/jwt"
)

var RsaPrivate string
var RsaPublic string

func ParseJwt(token string) (Jwt, error) {
	_, publicKey := getRsaKeys()
	encoder := jwtPkg.NewRsaEncoder(publicKey)
	claims, err := encoder.Decode(token)
	if err != nil {
		return nil, err
	}

	return &jwt{
		userId:            claims.Uint("user_id"),
		workspaceId:       claims.Uint("workspace_id"),
		workspaceMemberId: claims.Uint("workspace_member_id"),
		isAdmin:           claims.Bool("is_admin"),
		expiresAt:         claims.Int64("expires_at"),
	}, nil
}

func GenerateJwt(
	userId uint,
	workspaceId uint,
	workspaceMemberId uint,
	isAdmin bool,
	expiresAt int64,
) (string, error) {
	privateKey, _ := getRsaKeys()
	encoder := jwtPkg.NewRsaEncoder(privateKey)

	claimsBuilder := jwtPkg.NewClaimsBuilder().
		ExpirationTime(expiresAt).
		IssuedAt(time.Now().Unix()).
		Key("user_id").Set(userId).
		Key("workspace_id").Set(workspaceId).
		Key("workspace_member_id").Set(workspaceMemberId).
		Key("is_admin").Set(isAdmin).
		Build()
	token, err := encoder.Encode(claimsBuilder)
	if err != nil {
		return "", err
	}
	return token, nil
}

func getRsaKeys() (privateKey string, publicKey string) {
	// First, try to get keys from variables
	if RsaPrivate != "" && RsaPublic != "" {
		return RsaPrivate, RsaPublic
	}

	// Ensure temp directory exists
	tempDir := "temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		// If we can't create the directory, return empty strings
		return "", ""
	}

	// If not available, try to read from files
	privateKeyFile := filepath.Join(tempDir, "rsa.priv")
	publicKeyFile := filepath.Join(tempDir, "rsa.pub")

	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err == nil {
		publicKeyBytes, err := os.ReadFile(publicKeyFile)
		if err == nil {
			privateKey = string(privateKeyBytes)
			publicKey = string(publicKeyBytes)

			// Update variables for future use
			RsaPrivate = privateKey
			RsaPublic = publicKey

			return privateKey, publicKey
		}
	}

	// If files don't exist, generate new key pairs
	generator := jwtPkg.NewRsaGenerator()
	privateKey, publicKey, err = generator.GenerateKeyPair()
	if err != nil {
		// Return empty strings if generation fails
		return "", ""
	}

	// Save to files
	os.WriteFile(privateKeyFile, []byte(privateKey), 0644) // Read/write for owner only
	os.WriteFile(publicKeyFile, []byte(publicKey), 0644)   // Read for all, write for owner

	// Update variables
	RsaPrivate = privateKey
	RsaPublic = publicKey

	return privateKey, publicKey
}
