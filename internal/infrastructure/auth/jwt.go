// Package auth wraps RSA-backed JWT issue/verify so the rest of the app
// never touches go-jwt directly.
package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims mirrors lexik/jwt-bundle's default payload so PHP-issued tokens
// remain forward-compatible. Username typically carries the user's
// email; Roles is the role list copied from the user record.
type Claims struct {
	UID      int      `json:"uid,omitempty"`
	Username string   `json:"username"`
	Email    string   `json:"email,omitempty"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Service signs and validates JWT tokens using an RS256 keypair loaded
// from disk. Passphrase, if non-empty, is used to decrypt the private
// key (PKCS#8 PEM with encryption header).
type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttl        time.Duration
}

// NewService loads the keypair from disk. The public key is mandatory
// (verify path); the private key is optional (verify-only deployments).
func NewService(privateKeyPath, publicKeyPath, passphrase string, ttl time.Duration) (*Service, error) {
	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	s := &Service{publicKey: pub, ttl: ttl}

	if privateKeyPath != "" {
		privBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key: %w", err)
		}
		var priv *rsa.PrivateKey
		if passphrase != "" {
			priv, err = jwt.ParseRSAPrivateKeyFromPEMWithPassword(privBytes, passphrase)
		} else {
			priv, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
		}
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		s.privateKey = priv
	}
	return s, nil
}

// Issue signs a token for the given user.
func (s *Service) Issue(uid int, username, email string, roles []string) (string, error) {
	if s.privateKey == nil {
		return "", errors.New("jwt: private key not configured")
	}
	now := time.Now()
	claims := Claims{
		UID:      uid,
		Username: username,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

// ErrInvalidToken is returned by Verify when the token fails any
// validation step (signature, expiry, algorithm).
var ErrInvalidToken = errors.New("invalid token")

// Verify parses and validates a raw token string.
func (s *Service) Verify(raw string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	var claims Claims
	tok, err := parser.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		return s.publicKey, nil
	})
	if err != nil || !tok.Valid {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	return &claims, nil
}
