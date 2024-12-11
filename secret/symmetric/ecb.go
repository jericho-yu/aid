package symmetric

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jericho-yu/aid/compression"
	"github.com/jericho-yu/aid/str"
)

type (
	Ecb    struct{}
	Source struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

// padPKCS7 pads the plaintext to be a multiple of the block size
func (Ecb) padPKCS7(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// unPadPKCS7 removes the padding from the decrypted text
func (Ecb) unPadPKCS7(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}

func (Ecb) unPadPKCS72(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, fmt.Errorf("invalid blockSize: %d", blockSize)
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize {
		return nil, fmt.Errorf("invalid unpadding: %d", unpadding)
	}

	if unpadding == 0 {
		return nil, errors.New("invalid unpadding: 0")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

// Encrypt encrypts plaintext using AES in ECB mode
func (Ecb) Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext = Ecb{}.padPKCS7(plaintext, blockSize)
	cipherText := make([]byte, len(plaintext))

	for start := 0; start < len(plaintext); start += blockSize {
		block.Encrypt(cipherText[start:start+blockSize], plaintext[start:start+blockSize])
	}

	return cipherText, nil
}

// Decrypt decrypts cipherText using AES in ECB mode
func (Ecb) Decrypt(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(cipherText)%blockSize != 0 {
		return nil, fmt.Errorf("cipherText is not a multiple of the block size")
	}

	plaintext := make([]byte, len(cipherText))

	for start := 0; start < len(cipherText); start += blockSize {
		block.Decrypt(plaintext[start:start+blockSize], cipherText[start:start+blockSize])
	}

	return Ecb{}.unPadPKCS72(plaintext, blockSize)
}

// encrypt
func encrypt(plaintext, key []byte) string {
	// step1: zip
	zipped, zipErr := compression.NewZlib().Compress(plaintext)
	if zipErr != nil {
		str.NewTerminalLog("[ECB] compressing data: %v").Error(zipErr)
	}
	str.NewTerminalLog("[ECB] zipped").Info()

	// step2: aes-ecb-encrypt
	encrypted, encryptErr := Ecb{}.Encrypt(zipped, key)
	if encryptErr != nil {
		str.NewTerminalLog("[ECB] encrypting data:%v").Error(encryptErr)
	}
	str.NewTerminalLog("[ECB] encrypted").Info()

	// step3: encode base64
	encodeBase64 := base64.StdEncoding.EncodeToString(encrypted)
	str.NewTerminalLog("[ECB] encode base64").Success()

	return encodeBase64
}

// decrypt
func decrypt(cipherText, key []byte) []byte {
	// step1: decode base64
	decodeBase64, decodeBase64Err := base64.StdEncoding.DecodeString(string(cipherText))
	if decodeBase64Err != nil {
		str.NewTerminalLog("[ECB] decode base64: %v").Error(decodeBase64Err)
	}
	str.NewTerminalLog("[ECB] decode base64").Info()

	// step2: aes-ecb-decrypt
	decrypted, decryptErr := Ecb{}.Decrypt(decodeBase64, key)
	if decryptErr != nil {
		str.NewTerminalLog("[ECB] decrypting data: %v").Error(decryptErr)
	}
	str.NewTerminalLog("[ECB] decrypted").Info()

	// step3: unzip
	unzipped, unzipErr := compression.NewZlib().Decompress(decrypted)
	if unzipErr != nil {
		str.NewTerminalLog("[ECB] unzipping data: %v").Error(unzipErr)
	}
	str.NewTerminalLog("[ECB] unzipped").Info()

	return unzipped
}

func (Ecb) Demo() {
	key, keyErr := base64.StdEncoding.DecodeString("87dwQRkoNFNoIcq1A+zFHA==")
	if keyErr != nil {
		str.NewTerminalLog("[ECB] key decode base64: %v").Error(keyErr)
	}
	key2, key2Err := base64.StdEncoding.DecodeString("tjp5OPIU1ETF5s33fsLWdA==")
	if key2Err != nil {
		str.NewTerminalLog("[ECB] key2 decode base64: %v").Error(key2Err)
	}

	switch plaintext, jsonErr := json.Marshal(Source{Username: "abc123", Password: "cba321"}); {
	case jsonErr != nil:
		str.NewTerminalLog("[ECB] json marshal: %v").Error(jsonErr)
	default:
		str.NewTerminalLog("[ECB] json marshal: %s").Info(plaintext)

		// encrypt
		encrypted := encrypt(plaintext, key)
		str.NewTerminalLog("[ECB] encrypted: %s").Info(encrypted)

		// decrypt
		decrypted := decrypt([]byte(encrypted), key2)
		// decrypted := decrypt([]byte("Yw0Fh8699WC0hvKCRcFinq9nDqLdoECXZ5ZFK3onuXtzR61QKEvAH4+7NI4xYsn1eLwFOhzf0eBHnXv1ZaeWdueOy+t/OpgMXxl64s2PqLDRE8z+z2mHUVpvb7/V/2cS"), key2)
		str.NewTerminalLog("[ECB] decrypted: %s").Info(string(decrypted))
	}
}
