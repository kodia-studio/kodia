package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakeHashesPassword tests that Make generates a valid bcrypt hash
func TestMakeHashesPassword(t *testing.T) {
	password := "super-secret-password-123!"

	hash, err := Make(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

// TestMakeConsistency tests that Make hashes are consistent and verifiable
func TestMakeConsistency(t *testing.T) {
	password := "test-password"
	hash1, _ := Make(password)
	hash2, _ := Make(password)

	// Different hashes due to bcrypt salt
	assert.NotEqual(t, hash1, hash2)

	// But both should verify against the password
	assert.True(t, Check(password, hash1))
	assert.True(t, Check(password, hash2))
}

// TestCheckValidPassword tests verification of correct password
func TestCheckValidPassword(t *testing.T) {
	password := "my-password"
	hash, _ := Make(password)

	result := Check(password, hash)

	assert.True(t, result)
}

// TestCheckInvalidPassword tests that wrong password fails
func TestCheckInvalidPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"
	hash, _ := Make(password)

	result := Check(wrongPassword, hash)

	assert.False(t, result)
}

// TestCheckEmptyPassword tests empty password verification
func TestCheckEmptyPassword(t *testing.T) {
	hash, _ := Make("some-password")

	result := Check("", hash)

	assert.False(t, result)
}

// TestCheckEmptyHash tests empty hash verification
func TestCheckEmptyHash(t *testing.T) {
	result := Check("password", "")

	assert.False(t, result)
}

// TestCheckBothEmpty tests both empty inputs
func TestCheckBothEmpty(t *testing.T) {
	result := Check("", "")

	assert.False(t, result)
}

// TestMakeWithSpecialCharacters tests hashing passwords with special chars
func TestMakeWithSpecialCharacters(t *testing.T) {
	passwords := []string{
		"p@ss!word#123",
		"pässwörd",
		"密码123",
		"password\x00with\x00nulls", // Should still work (bcrypt doesn't stop at null)
	}

	for _, pwd := range passwords {
		hash, err := Make(pwd)
		require.NoError(t, err, "should hash password: %s", pwd)
		assert.True(t, Check(pwd, hash), "should verify password: %s", pwd)
	}
}

// TestMakeWithLongPassword tests hashing very long passwords
func TestMakeWithLongPassword(t *testing.T) {
	longPassword := "a" // 1000 character password
	for i := 0; i < 1000; i++ {
		longPassword += "a"
	}

	hash, err := Make(longPassword)

	require.NoError(t, err)
	assert.True(t, Check(longPassword, hash))
}

// TestCheckWithMalformedHash tests verification against malformed hash
func TestCheckWithMalformedHash(t *testing.T) {
	result := Check("password", "$2a$10$malformedhash")

	assert.False(t, result)
}

// BenchmarkMake benchmarks the Make function
func BenchmarkMake(b *testing.B) {
	password := "benchmark-password"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Make(password)
	}
}

// BenchmarkCheck benchmarks the Check function
func BenchmarkCheck(b *testing.B) {
	password := "benchmark-password"
	hash, _ := Make(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Check(password, hash)
	}
}
