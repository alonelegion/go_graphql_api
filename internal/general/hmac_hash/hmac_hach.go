package hmac_hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC interface {
	Hash(input string) string
}

type hm struct {
	hmac hash.Hash
}

// Функция создает экземпляры HMAC
func NewHMAC(key string) hm {
	h := hmac.New(sha256.New, []byte(key))
	return hm{
		hmac: h,
	}
}

// Функция будет хэшировать входную строку с секретным ключом
// который поставляется при создании объекта HMAC
func (h hm) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	hashedData := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedData)
}
