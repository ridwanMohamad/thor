package crypto_util

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

func ComputeHmac(value string, key string) string {
	secret := []byte(key)
	hm := hmac.New(sha512.New, secret)
	hm.Write([]byte(value))

	return hex.EncodeToString(hm.Sum(nil))
}
