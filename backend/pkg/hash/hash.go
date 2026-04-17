// Package hash provides password hashing and verification utilities.
// Uses bcrypt by default — NIST recommended for password storage.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

// Make hashes a plain-text password using bcrypt.
func Make(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Check verifies a plain-text password against a bcrypt hash.
// Returns true if the password matches.
func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
