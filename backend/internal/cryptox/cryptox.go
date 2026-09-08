// Package cryptox 收敛两类敏感数据的处理约定：
//
//   - 密码：bcrypt 哈希（不可逆）。永远不"加密"密码——能解出来的"加密"
//     意味着库被拖走时密码也一起泄密。bcrypt 自带盐值和慢哈希，防彩虹表。
//   - API Key：AES-256-GCM（可逆）。它必须能解出来（后端要替用户调 LLM），
//     所以用对称加密而不是哈希；GCM 同时提供机密性和完整性（防篡改）。
//
// 两把钥匙都来自服务端配置，绝不落库；库被拖走时没有主密钥拿不到明文。
package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// —— 密码：bcrypt ——

func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(h), nil
}

// CheckPassword 恒定时间比较由 bcrypt 内部保证，匹配返回 nil。
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// —— API Key：AES-256-GCM ——

// Box 持有主密钥派生的 AEAD。零值不可用，必须用 NewBox 构造。
type Box struct {
	aead cipher.AEAD
}

// NewBox 用 base64 编码的 32 字节主密钥构造。密钥长度错误立即报错——
// 这类配置错误要在启动时暴露，不能等第一次加密才发现。
func NewBox(keyB64 string) (*Box, error) {
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("aes key 不是合法 base64: %w", err)
	}
	if len(key) != 32 {
		return nil, errors.New("aes key 解码后必须是 32 字节（openssl rand -base64 32）")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Box{aead: aead}, nil
}

// Seal 加密明文，返回 base64(nonce + ciphertext)。每次调用随机 nonce，
// 同一明文两次加密结果不同（GCM 的语义安全性要求 nonce 不复用）。
func (b *Box) Seal(plaintext string) (string, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("rand nonce: %w", err)
	}
	ct := b.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

// Open 解密 Seal 的产物。密文被篡改或密钥不匹配都会报错（GCM 认证标签校验）。
func (b *Box) Open(encoded string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("密文不是合法 base64: %w", err)
	}
	ns := b.aead.NonceSize()
	if len(raw) < ns {
		return "", errors.New("密文太短")
	}
	pt, err := b.aead.Open(nil, raw[:ns], raw[ns:], nil)
	if err != nil {
		return "", errors.New("解密失败：密文损坏或主密钥不匹配")
	}
	return string(pt), nil
}

// ConstantTimeEqual 恒定时间字符串比较（验证码比对用，防时序侧信道）。
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
