package music

import (
	"MyChat/config"
	"MyChat/utils/file"
	"fmt"
	"log"
)

func IsExistMusicFile(user_id int64, file_id string) bool {
	pattern := fmt.Sprintf(config.GetConfig().MusicFilePath+"/%d/%s*", user_id, file_id)
	log.Println("file_path is ", pattern)
	return file.IsExistFile(pattern)
}
