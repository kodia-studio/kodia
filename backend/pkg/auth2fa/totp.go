// Package auth2fa provides TOTP based two-factor authentication utilities.
package auth2fa

import (
	"bytes"
	"image/png"
	"encoding/base64"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateSecret creates a new TOTP secret for a user.
// Returns the secret (base32) and the provisioning URL for QR code.
func GenerateSecret(userEmail string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Kodia",
		AccountName: userEmail,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// GenerateQRCodeBase64 returns a base64 encoded PNG of the QR code for a secret.
func GenerateQRCodeBase64(provisioningURL string) (string, error) {
	key, err := otp.NewKeyFromURL(provisioningURL)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}
	
	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// ValidateToken verifies a 6-digit TOTP token against a secret.
func ValidateToken(token, secret string) bool {
	return totp.Validate(token, secret)
}
