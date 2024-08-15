package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHmacSignature(payload []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return expectedSignature
}
