package music

import (
	"MyChat/common/mysql"
	"MyChat/config"
	"MyChat/model"
	"MyChat/utils/file"
	"fmt"
	"log"
)

func IsExistMusicFile(user_id int64, file_id string) bool {
	pattern := fmt.Sprintf(config.GetConfig().MusicFilePath+"/%d/%s*", user_id, file_id)
	log.Println("file_path is ", pattern)
	return file.IsExistFile(pattern)
}

func UploadMusicFile(uuid, music_name, file_path string, user_id, file_size int64, isupload int64) (*model.MusicFile, bool) {
	if music_file, err := mysql.InsertMusicFile(&model.MusicFile{
		UUID:      uuid,
		UserID:    user_id,
		MusicName: music_name,
		FilePath:  file_path,
		IsUpload:  isupload,
		FileSize:  file_size,
	}); err != nil {
		log.Println(err.Error())
		return nil, false
	} else {
		return music_file, true
	}

}
