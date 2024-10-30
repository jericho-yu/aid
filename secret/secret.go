package secret

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"

	"github.com/jericho-yu/aid/compression"

	"github.com/jericho-yu/aid/common"
	"github.com/jericho-yu/aid/secret/symmetric"
	"github.com/jericho-yu/aid/str"
)

type Secret struct{}

var SecretHelper Secret

func (Secret) New() *Secret { return &Secret{} }

func (Secret) EncryptAuthorization(key, secretKey string, iv []byte, randStr ...string) (string, string, error) {
	var (
		err   error
		uuid  = ""
		token []byte
	)

	if key == "" {
		return "", "", err
	}
	// 生成随机串
	if len(randStr) > 0 {
		uuid = randStr[0]
	} else {
		uuid, err = SecretHelper.MustEncrypt(str.NewRand().GetLetters(10))
		if err != nil {
			return "", "", err
		}
	}

	token, err = symmetric.CbcHelper.New().Encrypt([]byte(key+uuid), []byte(secretKey), iv)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(token), uuid, nil
}

func (Secret) DecryptAuthorization(token, secretKey string, iv []byte) (string, string, error) {
	var (
		err                   error
		token64, decryptToken []byte
	)

	if token == "" {
		return "", "", errors.New("token 不能为空")
	}

	token64, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", fmt.Errorf("base64解码token失败：%s", err.Error())
	}
	decryptToken, err = symmetric.CbcHelper.New().Decrypt(token64, []byte(secretKey), iv)
	if err != nil {
		return "", "", fmt.Errorf("解密失败：%s", err.Error())
	}

	return string(decryptToken[:len(decryptToken)-32]), string(decryptToken[len(decryptToken)-32:]), nil
}

func (Secret) MustEncrypt(data any) (string, error) {
	var (
		err       error
		dataBytes []byte
		h         hash.Hash
	)
	dataBytes = common.ToBytes(data)

	h = md5.New()
	if _, err = h.Write(dataBytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (Secret) Ecb16Encrypt(data any, needEncrypt, needZip bool, aes *symmetric.Aes) (string, error) {
	var (
		jsonByte, b                        []byte
		jsonMarshalErr, zipErr, encryptErr error
	)

	// json序列化
	jsonByte, jsonMarshalErr = json.Marshal(data)
	if jsonMarshalErr != nil {
		return "", jsonMarshalErr
	}

	// 压缩
	if needZip {
		b, zipErr = compression.NewZlib().Compress(jsonByte)
		if zipErr != nil {
			return "", zipErr
		}
	}

	// 加密
	if needEncrypt {
		b, encryptErr = symmetric.EcbHelper.New().Encrypt(b, aes.Encrypt.GetAesKey())
		if encryptErr != nil {
			return "", encryptErr
		}
	}

	if !needZip && !needEncrypt {
		return string(b), nil
	} else {
		return base64.StdEncoding.EncodeToString(b), nil
	}
}

func (Secret) Ecb16Decrypt(data string, needEncrypt, needZip bool, aes *symmetric.Aes) (any, error) {
	var (
		r                                                     any
		cipherText, decryptedByte, decompressedByte           []byte
		base64DecodeErr, jsonUnmarshalErr, decryptErr, zipErr error
	)

	if needEncrypt {
		// base64 解码
		cipherText, base64DecodeErr = base64.StdEncoding.DecodeString(data)
		if base64DecodeErr != nil {
			return nil, base64DecodeErr
		}

		// aes解密：ecb
		decryptedByte, decryptErr = symmetric.EcbHelper.New().Decrypt(cipherText, aes.Encrypt.GetAesKey())
		if decryptErr != nil {
			return nil, decryptErr
		}

		// 解压
		if needZip {
			decompressedByte, zipErr = compression.NewZlib().Decompress(decryptedByte)
			if zipErr != nil {
				return nil, zipErr
			}

			jsonUnmarshalErr = json.Unmarshal(decompressedByte, &r)
			if jsonUnmarshalErr != nil {
				return nil, jsonUnmarshalErr
			}

			return r, nil
		} else {
			// 将data反序列化
			jsonUnmarshalErr = json.Unmarshal(decryptedByte, &r)
			if jsonUnmarshalErr != nil {
				return nil, jsonUnmarshalErr
			}

			return r, nil
		}
	} else {
		jsonUnmarshalErr = json.Unmarshal([]byte(data), &r)
		if jsonUnmarshalErr != nil {
			return nil, jsonUnmarshalErr
		}
	}

	return r, nil
}
