package digest

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256 摘要算法
func Sha256(original []byte) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write(original); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)

	shaString := hex.EncodeToString(sum)

	return shaString, nil
}
