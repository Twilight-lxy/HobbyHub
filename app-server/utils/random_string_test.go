package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateRandomString 测试随机字符串生成函数
func TestGenerateRandomString(t *testing.T) {
	length := 16
	// 多次生成随机字符串，确保每次都不同
	var newRandomStringList []string
	for i := 0; i < 10; i++ {
		temp := GenerateRandomString(length)
		newRandomStringList = append(newRandomStringList, temp)
		// 检查生成的字符串长度
		assert.Equal(t, length, len(temp), "Generated string should have the correct length")
		// 检查生成的字符串是否只包含字母和数字
		for _, char := range temp {
			assert.True(t, (char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z'),
				"Generated string should only contain alphanumeric characters")
		}
	}
	// 检查newRandomStringList中所有字符串是否各不相同
	uniqueStrings := make(map[string]bool)
	for _, str := range newRandomStringList {
		uniqueStrings[str] = true
	}
	assert.Equal(t, len(newRandomStringList), len(uniqueStrings), "All generated random strings should be unique")
}
