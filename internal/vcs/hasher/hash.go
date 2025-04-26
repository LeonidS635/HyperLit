package hasher

import (
	"crypto/sha256"
	"fmt"
)

const HashSize = sha256.Size

func Calculate(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}

func ConvertToHex(hash []byte) string {
	return fmt.Sprintf("%x", hash)
}
