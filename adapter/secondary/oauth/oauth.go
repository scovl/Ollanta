// Package oauth provides OAuth provider implementations for GitHub, GitLab, and Google.
// Each provider implements domain/port.IOAuthProvider.
package oauth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateState returns a random hex string for CSRF protection.
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
