package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"regexp"
	"time"
)

func CheckMail(mail string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	return reg.MatchString(mail)
}
func CheckMobile(mobile string) bool {
	reg := regexp.MustCompile(`^1\d{10}$`)
	return reg.MatchString(mobile)
}
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func RandNumber(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

func RandString(l int, numberOnly bool) string {
	bytesTmp := []byte("0123456789")
	if !numberOnly {
		bytesTmp = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}

	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytesTmp[r.Intn(len(bytesTmp))])
	}
	return string(result)
}
