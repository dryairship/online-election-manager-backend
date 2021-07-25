package utils

import "github.com/mojocn/base64Captcha"

var store = base64Captcha.DefaultMemStore
var driver = base64Captcha.DefaultDriverDigit
var captcha = base64Captcha.NewCaptcha(driver, store)

func CreateCaptcha() (string, string, error) {
	return captcha.Generate()
}

func VerifyCaptcha(id, value string) bool {
	return store.Verify(id, value, true)
}
