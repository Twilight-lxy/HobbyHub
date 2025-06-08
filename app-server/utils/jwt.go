package utils

import (
	"hobbyhub-server/config"
	"hobbyhub-server/controllers"
	"hobbyhub-server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  u.Id,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // 72小时后过期
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().Authentication.JwtSecret))
}

func ParseJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们预期的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.GetConfig().Authentication.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	user := &models.User{
		Id: int64(claims["id"].(float64)),
	}

	user, err = controllers.GetUserByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
