package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	// ErrTokenExpired token已过期
	ErrTokenExpired = errors.New("token has expired")
	// ErrTokenInvalid token无效
	ErrTokenInvalid = errors.New("token is invalid")
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

const TokenExpireDuration = time.Hour * 24

var secretKey = []byte("your-secret-key") // TODO: 从配置文件读取

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
