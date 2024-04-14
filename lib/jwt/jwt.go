package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成token
func GenToken(username string) (string, error) {
	expires := viper.GetInt("expires")
	claims := CustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := viper.GetString("secret")

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(secret))
}

// ParseToken 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	secret := viper.GetString("secret")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
