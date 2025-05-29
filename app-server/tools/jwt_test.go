package tools

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// mock config for testing
func mockJwtSecret(secret string) func() {
	orig := config.GetConfig().Authentication.JwtSecret
	config.GetConfig().Authentication.JwtSecret = secret
	return func() {
		config.GetConfig().Authentication.JwtSecret = orig
	}
}

// 测试生成和解析JWT的正常流程
func TestGenerateJWTAndParseJWT(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	user := &models.User{ID: 123}
	tokenString, err := GenerateJWT(user)
	assert.NoError(t, err)          // 断言生成token无错误
	assert.NotEmpty(t, tokenString) // 断言token字符串不为空

	parsedUser, err := ParseJWT(tokenString)
	assert.NoError(t, err)                  // 断言解析token无错误
	assert.Equal(t, user.ID, parsedUser.ID) // 断言解析出的用户ID一致
}

// 测试解析非法token字符串
func TestParseJWT_InvalidToken(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	_, err := ParseJWT("invalid.token.string")
	assert.Error(t, err) // 断言解析非法token时返回错误
}

// 测试解析过期token
func TestParseJWT_ExpiredToken(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	claims := jwt.MapClaims{
		"id":  456,
		"exp": time.Now().Add(-time.Hour).Unix(), // 已过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("testsecret"))
	assert.NoError(t, err)

	_, err = ParseJWT(tokenString)
	assert.Error(t, err) // 断言解析过期token时返回错误
}

// 测试token签名不正确
func TestParseJWT_InvalidSignature(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	claims := jwt.MapClaims{
		"id":  789,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrongsecret")) // 用错误密钥签名
	assert.NoError(t, err)

	_, err = ParseJWT(tokenString)
	assert.Error(t, err) // 断言签名错误时返回错误
}
