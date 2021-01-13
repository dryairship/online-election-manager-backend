package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	mrand "math/rand"
)

const mPrefix = "GE"
const mSuffix = "2021"
const mPrefixLen = len(mPrefix)
const mSuffixLen = len(mSuffix)
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(n int32) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mrand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func get32SizeKey(key string) []byte {
	ans := make([]byte, 32)
	n := len(key)
	k := []byte(key)
	for i := range ans {
		ans[i] = k[i%n]
	}
	return ans
}

func EncryptAuthCode(key, message string) (string, error) {
	k := get32SizeKey(key)
	m := []byte(message)

	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(m)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return string(ciphertext), nil
}

func IsAuthCodeValid(key, encryptedText string) bool {
	text, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return false
	}
	k := get32SizeKey(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return false
	}
	if len(text) < aes.BlockSize {
		return false
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	str, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return false
	}
	bstr := []byte(str)
	n := len(bstr)
	return n > mPrefixLen+mSuffixLen && string(bstr[:mPrefixLen]) == mPrefix && string(bstr[n-mSuffixLen:]) == mSuffix
}

func GenerateAuthCode() string {
	return generateRandomString(12)
}

func GenerateMessage(roll string) string {
	return mPrefix + generateRandomString(12) + mSuffix
}
