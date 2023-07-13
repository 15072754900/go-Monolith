package utils

import (
	"errors"
	"fmt"
	"gin-blog-hufeng/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// 定义 token 相关 error
var (
	ErrTokenExpired     = errors.New("token 已过期，请重新登录")
	ErrTokenNotValidYet = errors.New("token 无效，请重新登录")
	ErrTokenMalformed   = errors.New("token 不正确，请重新登录")
	ErrTokenInvalid     = errors.New("这不是一个 token，请重新登录")
)

// MyClaims 定义 JWT 中的存储信息
type MyClaims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	UUID   string `json:"uuid"`
	jwt.RegisteredClaims
}

type MyJWT struct {
	Secret []byte
}

// GetJWT JWT 工具类 转换格式作用
func GetJWT() *MyJWT {
	return &MyJWT{[]byte(config.Cfg.JWT.Secret)}
}

// GetToken 生成 JWT
func (j *MyJWT) GetToken(userId int, role string, uuid string) (string, error) {
	claims := MyClaims{
		UserId: userId,
		Role:   role,
		UUID:   uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Cfg.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Cfg.JWT.Expire) * time.Hour)),
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的 secret 签名并获得完整编码后的字符串 token
	return token.SignedString(j.Secret)
}

// ParseToken 解析 JWT
func (j *MyJWT) ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.Secret, nil
		// 这个就是直接使用不需要进行过多的操作，但是一个很烦的现象，都是需要进行对外部的模块的二次包装，然后再使用
	})

	if err != nil {
		if vError, ok := err.(*jwt.ValidationError); ok {
			if vError.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if vError.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenExpired
			} else if vError.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else if vError.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}

	fmt.Println("正在校验")
	// 校验 token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	fmt.Println("未通过校验")

	return nil, ErrTokenInvalid
}
