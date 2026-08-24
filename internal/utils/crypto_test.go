package utils

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestEncryptDecryptSecret(t *testing.T) {
	old := viper.GetString("security.secret")
	viper.Set("security.secret", "test-secret-0123456789")
	defer viper.Set("security.secret", old)

	plain := "s3cr3t-password"
	enc, err := EncryptSecret(plain)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	if !strings.HasPrefix(enc, encryptedSecretPrefix) {
		t.Fatalf("密文缺少前缀: %q", enc)
	}
	if enc == plain {
		t.Fatal("密文不应与明文相同")
	}

	dec, err := DecryptSecret(enc)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	if dec != plain {
		t.Fatalf("解密结果不符: 期望 %q, 实际 %q", plain, dec)
	}
}

func TestEncryptSecretIdempotent(t *testing.T) {
	viper.Set("security.secret", "k1")
	enc, err := EncryptSecret("pw")
	if err != nil {
		t.Fatal(err)
	}
	enc2, err := EncryptSecret(enc)
	if err != nil || enc2 != enc {
		t.Fatalf("重复加密应返回原密文: %q vs %q, err=%v", enc, enc2, err)
	}
}

func TestDecryptSecretLegacyPlaintext(t *testing.T) {
	plain := "legacy-plain"
	dec, err := DecryptSecret(plain)
	if err != nil || dec != plain {
		t.Fatalf("历史明文应原样返回: %q, err=%v", dec, err)
	}
}

func TestEncryptSecretEmpty(t *testing.T) {
	enc, err := EncryptSecret("")
	if err != nil || enc != "" {
		t.Fatalf("空值应原样返回: %q, err=%v", enc, err)
	}
}

func TestEncryptSecretMissingKey(t *testing.T) {
	oldSec := viper.GetString("security.secret")
	oldJwt := viper.GetString("jwt.secret")
	viper.Set("security.secret", "")
	viper.Set("jwt.secret", "")
	defer func() {
		viper.Set("security.secret", oldSec)
		viper.Set("jwt.secret", oldJwt)
	}()

	enc, err := EncryptSecret("pw")
	if err == nil {
		t.Fatal("缺少密钥时应返回错误")
	}
	if enc != "pw" {
		t.Fatalf("缺少密钥时应返回原值, 实际 %q", enc)
	}
}

func TestDecryptSecretTampered(t *testing.T) {
	viper.Set("security.secret", "k2")
	enc, err := EncryptSecret("pw")
	if err != nil {
		t.Fatal(err)
	}
	tampered := enc[:len(enc)-2] + "XX"
	if _, err := DecryptSecret(tampered); err == nil {
		t.Fatal("篡改的密文应解密失败")
	}
}
