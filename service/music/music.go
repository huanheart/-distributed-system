package music

import (
	"MyChat/common/mysql"
	"MyChat/dao/music"
	"MyChat/model"
)

func IsExistMusicFile(user_id int64, file_id string) (*model.MusicFile, bool) {

	//1:先查数据库是否存在这个user_id和file_id都含有的文件
	musicfile, err := mysql.GetMusicfile(user_id, file_id)
	if err != nil {
		return nil, false
	}
	//todo:后续可以将本地存储存储到oss上，弄一个抽象类,这样就不需要真实在存储到本地，可在oss与本地存储等任选
	//2:查看是否含有/static/user_id/file_id为前缀的文件
	if ok := music.IsExistMusicFile(user_id, file_id); !ok {
		return nil, false
	}

	return musicfile, true
}
