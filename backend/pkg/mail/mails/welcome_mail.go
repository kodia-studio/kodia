package mails

import (
	"context"
	"fmt"

	"github.com/kodia-studio/kodia/pkg/mail"
)

// WelcomeMail is sent to new users upon registration
type WelcomeMail struct {
	*mail.TemplateMail
	UserName string
	Email    string
}

// NewWelcomeMail creates a new welcome mail
func NewWelcomeMail(templateEngine *mail.TemplateEngine) *WelcomeMail {
	return &WelcomeMail{
		TemplateMail: mail.NewTemplateMail(templateEngine),
	}
}

// To sets the recipient
func (wm *WelcomeMail) SetEmail(email string) *WelcomeMail {
	wm.Email = email
	wm.To(email)
	return wm
}

// SetUserName sets the user name
func (wm *WelcomeMail) SetUserName(name string) *WelcomeMail {
	wm.UserName = name
	wm.SetTemplateVariable("name", name)
	return wm
}

// Build constructs the welcome email
func (wm *WelcomeMail) Build(ctx context.Context, data map[string]any) error {
	// Set subject
	wm.Subject(fmt.Sprintf("Welcome to Kodia, %s!", wm.UserName))

	// Set template name
	wm.WithTemplate("welcome")

	// Add data to template
	wm.SetTemplateVariable("name", wm.UserName)
	wm.SetTemplateVariable("email", wm.Email)
	wm.SetTemplateVariable("app_name", "Kodia")

	// Render template
	return wm.TemplateMail.Build(ctx, data)
}

// VerificationMail is sent to verify user email
type VerificationMail struct {
	*mail.TemplateMail
	UserName    string
	Email       string
	VerifyLink  string
	ExpiresIn   int // minutes
}

// NewVerificationMail creates a new verification mail
func NewVerificationMail(templateEngine *mail.TemplateEngine) *VerificationMail {
	return &VerificationMail{
		TemplateMail: mail.NewTemplateMail(templateEngine),
		ExpiresIn:    24, // 24 hours default
	}
}

// SetEmail sets the recipient email
func (vm *VerificationMail) SetEmail(email string) *VerificationMail {
	vm.Email = email
	vm.To(email)
	return vm
}

// SetUserName sets the user name
func (vm *VerificationMail) SetUserName(name string) *VerificationMail {
	vm.UserName = name
	return vm
}

// SetVerifyLink sets the verification link
func (vm *VerificationMail) SetVerifyLink(link string) *VerificationMail {
	vm.VerifyLink = link
	return vm
}

// SetExpiresIn sets the expiration time in minutes
func (vm *VerificationMail) SetExpiresIn(minutes int) *VerificationMail {
	vm.ExpiresIn = minutes
	return vm
}

// Build constructs the verification email
func (vm *VerificationMail) Build(ctx context.Context, data map[string]any) error {
	// Set subject
	vm.Subject("Verify Your Email Address")

	// Set template name
	vm.WithTemplate("verification")

	// Add data to template
	vm.SetTemplateVariable("name", vm.UserName)
	vm.SetTemplateVariable("email", vm.Email)
	vm.SetTemplateVariable("verify_link", vm.VerifyLink)
	vm.SetTemplateVariable("expires_in", vm.ExpiresIn)
	vm.SetTemplateVariable("app_name", "Kodia")

	// Render template
	return vm.TemplateMail.Build(ctx, data)
}

// PasswordResetMail is sent for password reset requests
type PasswordResetMail struct {
	*mail.TemplateMail
	UserName   string
	Email      string
	ResetLink  string
	ExpiresIn  int // minutes
}

// NewPasswordResetMail creates a new password reset mail
func NewPasswordResetMail(templateEngine *mail.TemplateEngine) *PasswordResetMail {
	return &PasswordResetMail{
		TemplateMail: mail.NewTemplateMail(templateEngine),
		ExpiresIn:    60, // 60 minutes default
	}
}

// SetEmail sets the recipient email
func (prm *PasswordResetMail) SetEmail(email string) *PasswordResetMail {
	prm.Email = email
	prm.To(email)
	return prm
}

// SetUserName sets the user name
func (prm *PasswordResetMail) SetUserName(name string) *PasswordResetMail {
	prm.UserName = name
	return prm
}

// SetResetLink sets the password reset link
func (prm *PasswordResetMail) SetResetLink(link string) *PasswordResetMail {
	prm.ResetLink = link
	return prm
}

// SetExpiresIn sets the expiration time in minutes
func (prm *PasswordResetMail) SetExpiresIn(minutes int) *PasswordResetMail {
	prm.ExpiresIn = minutes
	return prm
}

// Build constructs the password reset email
func (prm *PasswordResetMail) Build(ctx context.Context, data map[string]any) error {
	// Set subject
	prm.Subject("Reset Your Password")

	// Set template name
	prm.WithTemplate("password-reset")

	// Add data to template
	prm.SetTemplateVariable("name", prm.UserName)
	prm.SetTemplateVariable("email", prm.Email)
	prm.SetTemplateVariable("reset_link", prm.ResetLink)
	prm.SetTemplateVariable("expires_in", prm.ExpiresIn)
	prm.SetTemplateVariable("app_name", "Kodia")

	// Render template
	return prm.TemplateMail.Build(ctx, data)
}
