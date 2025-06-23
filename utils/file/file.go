package file

import (
	"log"
	"path/filepath"
)

// 查找一个文件是否存在
func IsExistFile(file_path string) bool {
	// 查找匹配的文件列表
	matches, err := filepath.Glob(file_path)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return len(matches) > 0
}
