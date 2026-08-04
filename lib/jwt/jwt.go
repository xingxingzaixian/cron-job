package jwt

import (
	"cronJob/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成token
func GenToken(username string) (string, error) {
	expires := viper.GetInt("jwt.expires")
	claims := CustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expires))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := signingSecret()

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(secret))
}

// ParseToken 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	secret := signingSecret()
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// signingSecret 获取JWT签名密钥
// 未配置或为空时使用随机密钥并告警，避免空密钥签名导致token可被伪造；
// 注意：随机密钥模式下服务重启后所有登录态失效，生产环境应固定配置 jwt.secret
func signingSecret() string {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = utils.GenerateRandomJWTSecret()
		zap.S().Warn("jwt.secret 未配置或为空，使用随机密钥签名，重启后所有登录态将失效")
	}
	return secret
}
