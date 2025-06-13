package utils

import (
	"hobbyhub-server/config"
	"hobbyhub-server/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mock config for testing
func mockJwtSecret(secret string) func() {
	orig := config.GetConfig().Authentication.JwtSecret
	config.GetConfig().Authentication.JwtSecret = secret
	return func() {
		config.GetConfig().Authentication.JwtSecret = orig
	}
}

// 添加到jwt_test.go文件中
func setupDBMock() func() {
	// 使用类似于之前看到的SetupMockDB方法
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// 保存原始DB
	origDB := config.DB
	// 设置mock DB
	config.DB = gormDB

	// 为模拟GetUserByUserId的调用准备数据
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE id = ? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs(int64(123), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
			AddRow(123, "test_user"))

	return func() {
		config.DB = origDB
	}
}

// 测试生成和解析JWT的正常流程
func TestGenerateJWTAndParseJWT(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	// 添加数据库mock
	dbRestore := setupDBMock()
	defer dbRestore()

	user := &models.User{Id: 123}
	tokenString, err := GenerateJWT(user)
	assert.NoError(t, err)          // 断言生成token无错误
	assert.NotEmpty(t, tokenString) // 断言token字符串不为空

	parsedUser, err := ParseJWT(tokenString)
	assert.NoError(t, err)                  // 断言解析token无错误
	assert.Equal(t, user.Id, parsedUser.Id) // 断言解析出的用户ID一致
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
		"id":   456,
		"role": "admin",
		"exp":  jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 已过期
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
		"id":   789,
		"role": "user",
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrongsecret")) // 用错误密钥签名
	assert.NoError(t, err)

	_, err = ParseJWT(tokenString)
	assert.Error(t, err) // 断言签名错误时返回错误
}

func mockJwtSecretAdmin(secret string) func() {
	orig := config.GetConfig().Authentication.JwtSecret
	config.GetConfig().Authentication.JwtSecret = secret
	return func() {
		config.GetConfig().Authentication.JwtSecret = orig
	}
}
func setupDBMockAdmin() func() {
	// 使用类似于之前看到的SetupMockDB方法
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// 保存原始DB
	origDB := config.DB
	// 设置mock DB
	config.DB = gormDB

	// 为模拟GetAdminById的调用准备数据
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `admin` WHERE id = ? ORDER BY `admin`.`id` LIMIT ?")).
		WithArgs(int64(123), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
			AddRow(123, "test_admin"))

	return func() {
		config.DB = origDB
	}
}

// 测试生成和解析管理员JWT的正常流程
func TestGenerateAdminJWTAndParseAdminJWT(t *testing.T) {
	restore := mockJwtSecretAdmin("testsecret")
	defer restore()

	// 添加数据库mock
	dbRestore := setupDBMockAdmin()
	defer dbRestore()

	admin := &models.Admin{Id: 123}
	tokenString, err := GenerateAdminJWT(admin)
	assert.NoError(t, err)          // 断言生成token无错误
	assert.NotEmpty(t, tokenString) // 断言token字符串不为空

	parsedAdmin, err := ParseAdminJWT(tokenString)
	assert.NoError(t, err)                    // 断言解析token无错误
	assert.Equal(t, admin.Id, parsedAdmin.Id) // 断言解析出的管理员ID一致
}

// 测试解析非法管理员token字符串
func TestParseAdminJWT_InvalidToken(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	_, err := ParseAdminJWT("invalid.token.string")
	assert.Error(t, err) // 断言解析非法token时返回错误
}

// 测试解析过期管理员token
func TestParseAdminJWT_ExpiredToken(t *testing.T) {
	restore := mockJwtSecret("testsecret")
	defer restore()

	claims := jwt.MapClaims{
		"id":   456,
		"role": "admin",
		"exp":  jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 已过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("testsecret"))
	assert.NoError(t, err)

	_, err = ParseAdminJWT(tokenString)
	assert.Error(t, err) // 断言解析过期token时返回错误
}
func TestParseAdminJWT_InvalidSignature(t *testing.T) {
	restore := mockJwtSecretAdmin("testsecret")
	defer restore()
	claims := jwt.MapClaims{
		"id":   789,
		"role": "admin",
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour)), // 未过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("wrongsecret")) // 用错误密钥签名
	assert.NoError(t, err)

	_, err = ParseAdminJWT(tokenString)
	assert.Error(t, err) // 断言签名错误时返回错误
}
func TestParseAdminJWT_InvalidRole(t *testing.T) {
	restore := mockJwtSecretAdmin("testsecret")
	defer restore()

	claims := jwt.MapClaims{
		"id":   789,
		"role": "user",
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour)), // 未过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("testsecret"))
	assert.NoError(t, err)

	_, err = ParseAdminJWT(tokenString)
	assert.Error(t, err) // 断言角色不正确时返回错误
}
