package utils

import (
	"crypto/rand"
)

func GenerateRandomString(length int) string {
	// 定义字符集，包含所有字母和数字
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length) // 创建指定长度的字节切片
	_, err := rand.Read(b)    // 使用crypto/rand填充随机字节
	if err != nil {
		panic(err)
	}
	for i := range b { // 遍历每个字节，将随机字节映射到字符集范围内的字符
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b) // 将字节切片转换为字符串并返回
}
