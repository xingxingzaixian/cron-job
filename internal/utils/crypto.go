package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// encryptedSecretPrefix 加密值前缀，用于区分历史明文与密文
const encryptedSecretPrefix = "enc:v1:"

// encryptionKey 从配置派生AES-256密钥
// 优先使用 security.secret，未配置时回退到 jwt.secret；两者都为空则返回错误
func encryptionKey() ([]byte, error) {
	secret := viper.GetString("security.secret")
	if secret == "" {
		secret = viper.GetString("jwt.secret")
	}
	if secret == "" {
		return nil, errors.New("未配置 security.secret/jwt.secret，无法加密敏感数据")
	}
	sum := sha256.Sum256([]byte(secret))
	return sum[:], nil
}

// EncryptSecret 加密敏感字符串（AES-256-GCM）
// 空值或已加密的值原样返回（幂等）；加密失败时返回原值并告警
func EncryptSecret(plaintext string) (string, error) {
	if plaintext == "" || strings.HasPrefix(plaintext, encryptedSecretPrefix) {
		return plaintext, nil
	}

	key, err := encryptionKey()
	if err != nil {
		zap.S().Warnf("敏感数据加密失败（将按明文存储）: %v", err)
		return plaintext, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext, err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)
	return encryptedSecretPrefix + base64.StdEncoding.EncodeToString(payload), nil
}

// DecryptSecret 解密敏感字符串；非加密值（历史明文）原样返回
func DecryptSecret(value string) (string, error) {
	if !strings.HasPrefix(value, encryptedSecretPrefix) {
		return value, nil
	}

	key, err := encryptionKey()
	if err != nil {
		return "", err
	}

	encoded := strings.TrimPrefix(value, encryptedSecretPrefix)
	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("密文解码失败: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(payload) < nonceSize {
		return "", errors.New("密文长度非法")
	}
	plaintext, err := gcm.Open(nil, payload[:nonceSize], payload[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("密文解密失败(密钥变更或数据损坏): %v", err)
	}
	return string(plaintext), nil
}
