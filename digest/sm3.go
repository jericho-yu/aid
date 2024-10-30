package digest

import (
	"github.com/tjfoc/gmsm/sm3"

	"encoding/hex"
)

// Sm3 生成sm3摘要
func Sm3(original []byte) string {
	h := sm3.New()
	h.Write(original)
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}
