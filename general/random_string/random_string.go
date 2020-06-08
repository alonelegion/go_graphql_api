package random_string

import (
	"crypto/rand"
	"encoding/base64"
)

const TokenBytes = 32

type RandomString interface {
	GenerateToken() (string, error)
	NumberOfBytes(base64String string) (int, error)
}

type randomString struct{}

func NewRandomString() RandomString {
	return &randomString{}
}

// Функция генерации токена с заданым размером байт
// Генерация среза байт размером n-байт, а затем возвращает строку
// в base64URL кодировке этого среза байт
func (r *randomString) GenerateToken() (string, error) {
	b, err := r.generateRandomBytes(TokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Функция возвращает количество байт в любой строке
func (r *randomString) NumberOfBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

// Функция генерации среза
func (r *randomString) generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
