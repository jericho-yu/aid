package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// hashString 生成字符串的 SHA-256 哈希值
func hashString(s string) int {
	hash := sha256.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)
	return int(binary.BigEndian.Uint32(hashBytes[:4]))
}

func main() {
	str := "AinfinitCabinetExtendedDevices-com.jz.autobooth!"
	hashedInt := hashString(str)
	hashedHex := fmt.Sprintf("%x", hashedInt)
	fmt.Println("Original String:", str)
	fmt.Println("Hashed Int:", hashedInt)
	fmt.Println("Hashed Hex:", hashedHex)
}
