package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidatorNew creates a new validator
func TestValidatorNew(t *testing.T) {
	v := New()

	require.NotNil(t, v)
	require.NotNil(t, v.Engine())
}

// TestValidatorStruct tests struct validation
func TestValidatorStruct(t *testing.T) {
	type User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	v := New()

	// Valid user
	validUser := User{Email: "test@example.com", Password: "password123"}
	err := v.Struct(validUser)
	assert.NoError(t, err)

	// Missing email
	invalidUser := User{Email: "", Password: "password123"}
	err = v.Struct(invalidUser)
	assert.Error(t, err)

	// Invalid email format
	invalidEmail := User{Email: "not-an-email", Password: "password123"}
	err = v.Struct(invalidEmail)
	assert.Error(t, err)

	// Password too short
	shortPassword := User{Email: "test@example.com", Password: "short"}
	err = v.Struct(shortPassword)
	assert.Error(t, err)
}

// TestFormatErrors tests error formatting
func TestFormatErrors(t *testing.T) {
	type TestStruct struct {
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"required,min=18"`
	}

	v := New()
	invalid := TestStruct{Email: "invalid-email", Age: 10}

	err := v.Struct(invalid)
	require.Error(t, err)

	formatted := FormatErrors(err)

	// Check that the formatted errors have the fields
	assert.True(t, len(formatted["email"]) > 0)
	assert.True(t, len(formatted["age"]) > 0)

	// Check that error messages contain field names
	emailErrorFound := false
	for _, msg := range formatted["email"] {
		if len(msg) > 0 {
			emailErrorFound = true
			break
		}
	}
	assert.True(t, emailErrorFound)

	ageErrorFound := false
	for _, msg := range formatted["age"] {
		if len(msg) > 0 {
			ageErrorFound = true
			break
		}
	}
	assert.True(t, ageErrorFound)
}

// TestValidationTags tests standard validation tags
func TestValidationTags(t *testing.T) {
	type TestData struct {
		Required string `json:"required_field" validate:"required"`
		Email    string `json:"email_field" validate:"required,email"`
		Min      string `json:"min_field" validate:"required,min=3"`
		Max      string `json:"max_field" validate:"required,max=5"`
		Len      string `json:"len_field" validate:"required,len=5"`
		URL      string `json:"url_field" validate:"required,url"`
	}

	v := New()

	tests := []struct {
		name    string
		data    TestData
		isValid bool
	}{
		{
			name: "Valid data",
			data: TestData{
				Required: "value",
				Email:    "test@example.com",
				Min:      "abc",
				Max:      "abcd",
				Len:      "exact",
				URL:      "https://example.com",
			},
			isValid: true,
		},
		{
			name: "Missing required field",
			data: TestData{
				Email: "test@example.com",
				Min:   "abc",
				Max:   "ab",
				Len:   "exact",
				URL:   "https://example.com",
			},
			isValid: false,
		},
		{
			name: "Invalid email",
			data: TestData{
				Required: "value",
				Email:    "not-an-email",
				Min:      "abc",
				Max:      "ab",
				Len:      "exact",
				URL:      "https://example.com",
			},
			isValid: false,
		},
		{
			name: "Min validation fails",
			data: TestData{
				Required: "value",
				Email:    "test@example.com",
				Min:      "ab", // Less than 3
				Max:      "ab",
				Len:      "exact",
				URL:      "https://example.com",
			},
			isValid: false,
		},
		{
			name: "Max validation fails",
			data: TestData{
				Required: "value",
				Email:    "test@example.com",
				Min:      "abc",
				Max:      "abcdef", // More than 5
				Len:      "exact",
				URL:      "https://example.com",
			},
			isValid: false,
		},
		{
			name: "Len validation fails",
			data: TestData{
				Required: "value",
				Email:    "test@example.com",
				Min:      "abc",
				Max:      "ab",
				Len:      "abc", // Not exactly 5
				URL:      "https://example.com",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.data)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestStrongPasswordValidation tests strong password rule
func TestStrongPasswordValidation(t *testing.T) {
	type Request struct {
		Password string `json:"password" validate:"required,strong_password"`
	}

	v := New()

	tests := []struct {
		name     string
		password string
		isValid  bool
	}{
		{"Valid strong password", "SecurePass123!", true},
		{"Missing uppercase", "securepass123!", false},
		{"Missing lowercase", "SECUREPASS123!", false},
		{"Missing digit", "SecurePassd!", false},
		{"Missing symbol", "SecurePass123", false},
		{"All requirements met", "P@ssw0rd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{Password: tt.password}
			err := v.Struct(req)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestPhoneValidation tests phone validation rule
func TestPhoneValidation(t *testing.T) {
	type Request struct {
		Phone string `json:"phone" validate:"required,phone"`
	}

	v := New()

	tests := []struct {
		name    string
		phone   string
		isValid bool
	}{
		{"Valid E.164 format", "+1234567890", true},
		{"With spaces", "+1 234 567 890", true},
		{"With dashes", "+1-234-567-890", true},
		{"With parentheses", "+1(234)567-890", true},
		{"Too short", "12345", false},
		{"Too long", "123456789012345678901", false},
		{"Valid 10 digits", "1234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{Phone: tt.phone}
			err := v.Struct(req)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestAlphaSpaceValidation tests alpha_space validation rule
func TestAlphaSpaceValidation(t *testing.T) {
	type Request struct {
		Name string `json:"name" validate:"required,alpha_space"`
	}

	v := New()

	tests := []struct {
		name    string
		value   string
		isValid bool
	}{
		{"Simple name", "John Doe", true},
		{"Single word", "John", true},
		{"Multiple words", "John Adam Smith", true},
		{"With numbers", "John123", false},
		{"With punctuation", "John-Doe", false},
		{"With special chars", "John@Doe", false},
		{"Empty string after required check", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{Name: tt.value}
			err := v.Struct(req)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestNoHTMLValidation tests no_html validation rule
func TestNoHTMLValidation(t *testing.T) {
	type Request struct {
		Comment string `json:"comment" validate:"required,no_html"`
	}

	v := New()

	tests := []struct {
		name    string
		value   string
		isValid bool
	}{
		{"Plain text", "This is a comment", true},
		{"With special chars", "Comment with @#$% chars", true},
		{"With HTML tag open", "Comment <script>", false},
		{"With HTML tag close", "Comment >alert", false},
		{"Full HTML", "<p>Comment</p>", false},
		{"Safe special chars", "Comment with & and quotes", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{Comment: tt.value}
			err := v.Struct(req)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestFormatErrorsFormatting tests FormatErrors output structure
func TestFormatErrorsFormatting(t *testing.T) {
	type Request struct {
		Email string `json:"email" validate:"required,email"`
		Name  string `json:"name" validate:"required,min=3"`
	}

	v := New()
	req := Request{Email: "invalid", Name: "ab"}

	err := v.Struct(req)
	require.Error(t, err)

	formatted := FormatErrors(err)

	// Should have entries for both fields
	assert.Contains(t, formatted, "email")
	assert.Contains(t, formatted, "name")

	// Each field should have error messages
	assert.True(t, len(formatted["email"]) > 0)
	assert.True(t, len(formatted["name"]) > 0)

	// Error messages should be strings
	for field, messages := range formatted {
		assert.True(t, len(messages) > 0, "field %s should have messages", field)
		for _, msg := range messages {
			assert.NotEmpty(t, msg)
		}
	}
}

// TestNestedStructValidation tests validation of nested structs
func TestNestedStructValidation(t *testing.T) {
	type Address struct {
		City string `json:"city" validate:"required"`
	}

	type User struct {
		Email   string  `json:"email" validate:"required,email"`
		Address Address `validate:"required"`
	}

	v := New()

	// Valid nested
	user := User{
		Email: "test@example.com",
		Address: Address{
			City: "New York",
		},
	}
	err := v.Struct(user)
	assert.NoError(t, err)

	// Invalid nested
	invalidUser := User{
		Email: "test@example.com",
		Address: Address{
			City: "", // Empty city
		},
	}
	err = v.Struct(invalidUser)
	assert.Error(t, err)
}

// BenchmarkValidate benchmarks struct validation
func BenchmarkValidate(b *testing.B) {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,strong_password"`
	}

	v := New()
	req := Request{Email: "test@example.com", Password: "SecurePass123!"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Struct(req)
	}
}

// BenchmarkFormatErrors benchmarks error formatting
func BenchmarkFormatErrors(b *testing.B) {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	v := New()
	req := Request{} // Invalid

	err := v.Struct(req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatErrors(err)
	}
}
