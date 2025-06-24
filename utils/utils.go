package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		// 0~9随机数
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func GenerateUUID() string {
	return uuid.New().String()
}

// 获取文件名
func GetFilePreName(filename, fileExt string) string {
	return strings.TrimSuffix(filename, fileExt) // "hello_song"
}
