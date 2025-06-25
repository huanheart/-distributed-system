package utils

import (
	"MyChat/config"
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

// file_path是真实路径，我们需要给他变成服务路径
// 当前是http（后续要在config.GetConfig().MusicFileIp加上http/https）
func GetHttpPath(file_path string) string {
	base := config.GetConfig().MusicFilePath // D:/MyMusicPlatform/test
	full := file_path                        // D:/MyMusicPlatform/test/1.jpg
	rest := strings.TrimPrefix(full, base)   // -> /1.jpg
	res := config.GetConfig().MusicFileIp + ":" + strconv.Itoa(config.GetConfig().Port) + config.GetConfig().HttpFilePath + rest
	return res
}
