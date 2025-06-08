package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHashPassword 测试密码哈希函数
func TestHashPassword(t *testing.T) {
	password := "testPassword"
	expectedHash := "fd5cb51bafd60f6fdbedde6e62c473da6f247db271633e15919bab78a02ee9eb" // SHA-256 hash of "testPassword"

	hashedPassword := HashPassword(password)
	assert.Equal(t, expectedHash, hashedPassword, "Hashed password should match expected hash")
}

// TestCheckPasswordHash 测试密码验证函数
func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword"
	expectedHash := HashPassword(password)

	// 测试正确密码
	assert.True(t, CheckPasswordHash(password, expectedHash), "Password should match the hash")

	// 测试错误密码
	assert.False(t, CheckPasswordHash("wrongPassword", expectedHash), "Wrong password should not match the hash")
}
