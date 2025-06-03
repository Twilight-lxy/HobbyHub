package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 使用SHA-256算法对密码进行哈希
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHash 验证密码是否与哈希值匹配
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
