package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHash(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	res := hex.EncodeToString(hash.Sum(nil))

	return res
}
