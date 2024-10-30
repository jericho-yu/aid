package asymmetric

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type (
	PemBase64 struct {
		base64PublicKey  string
		base64PrivateKey string
		publicKey        []byte
		privateKey       []byte
	}
)

var PemBase64Helper PemBase64

func (PemBase64) New() *PemBase64 { return &PemBase64Helper }

func (r *PemBase64) SetBase64PublicKey(base64PublicKey string) *PemBase64 {
	r.base64PublicKey = base64PublicKey
	return r
}

func (r *PemBase64) SetBase64PrivateKye(base64PrivateKye string) *PemBase64 {
	r.base64PrivateKey = base64PrivateKye
	return r
}

func (r *PemBase64) GetBase64PublicKey() string {
	return r.base64PublicKey
}

func (r *PemBase64) GetBase64PrivateKey() string {
	return r.base64PrivateKey
}

func (r *PemBase64) GetPemPublicKey() []byte {
	return r.publicKey
}

func (r *PemBase64) GeneratePemPublicKey() (*PemBase64, error) {
	// 解码Base64字符串
	publicKeyBytes, err := base64.StdEncoding.DecodeString(r.base64PublicKey)
	if err != nil {
		return r, fmt.Errorf("解码Base64失败: %v", err)
	}

	// 尝试解析为PEM块
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		// 如果不是PEM格式，则尝试解析为x509公钥并创建一个PEM块
		_, err = x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return r, fmt.Errorf("解析公钥失败se64失败: %v", err)
		}

		// 创建PEM块
		block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
	}

	// 将PEM块编码为内存中的字节切片
	r.publicKey = pem.EncodeToMemory(block)

	return r, nil
}

// GetPemPrivateKey 获取pem私钥
func (r *PemBase64) GetPemPrivateKey() []byte {
	return r.privateKey
}

// GeneratePemPrivateKey 生成pem密钥
func (r *PemBase64) GeneratePemPrivateKey() (*PemBase64, error) {
	// 解码Base64字符串
	privateKeyBytes, err := base64.StdEncoding.DecodeString(r.base64PrivateKey)
	if err != nil {
		return r, fmt.Errorf("解码Base64失败: %v", err)
	}

	// 手动添加PEM头部和尾部
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	r.privateKey = pem.EncodeToMemory(pemBlock)

	// 尝试解析为PEM块
	block, _ := pem.Decode(r.privateKey)
	if block == nil {
		return r, errors.New("不是有效的PEM编码私钥")
	}

	return r, nil
}
