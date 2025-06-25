package music

import (
	"MyChat/common/mysql"
	"MyChat/common/rabbitmq"
	"MyChat/config"
	"MyChat/dao/music"
	"MyChat/model"
	"MyChat/utils"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func IsExistMusicFile(user_id int64, file_id string) (*model.MusicFile, bool) {

	//1:先查数据库是否存在这个user_id和file_id都含有的文件
	musicfile, err := mysql.GetMusicfile(user_id, file_id)
	if err != nil {
		return nil, false
	}
	//todo:后续可以将本地存储存储到oss上，弄一个抽象类,这样就不需要真实在存储到本地，可在oss与本地存储等任选
	//2:查看是否含有/static/user_id/file_id为前缀的文件
	//todo:这个也需要弄一个抽象，因为如果上传到本地和远端的检查方式是不同的
	if ok := music.IsExistMusicFile(user_id, file_id); !ok {
		return nil, false
	}
	//说明此时还没有上传完毕
	if musicfile.IsUpload == 0 {
		return nil, false
	}

	return musicfile, true
}

func MusicUpload(user_id int64, file *multipart.FileHeader) (*model.MusicFile, bool) {
	// 1. 生成 UUID
	fileID := utils.GenerateUUID()

	// 2. 提取扩展名
	fileExt := filepath.Ext(file.Filename)
	saveName := fileID + fileExt
	savePath := fmt.Sprintf("%s/%d/%s", config.GetConfig().MusicFilePath, user_id, saveName)

	// 3. 创建目录
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		log.Println(err.Error())
		return nil, false
	}

	// 4. 保存文件
	src, err := file.Open()
	if err != nil {
		log.Println(err.Error())
		return nil, false
	}
	defer src.Close()

	out, err := os.Create(savePath)
	if err != nil {
		log.Println(err.Error())
		return nil, false
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, false
	}
	music_file, ok := music.UploadMusicFile(fileID, utils.GetFilePreName(file.Filename, fileExt), savePath, user_id, file.Size, 0)
	if !ok {
		return nil, ok
	}
	// 7. 将推送消息队列任务（异步上传 OSS） 后续封装了消息队列部分改这个部分
	//todo:后续弄一层抽象，弄一个抽象父类，对应方法都一样，然后只需要在外层写new 哪一个子类即可 ,代码维护量大大减少
	data := rabbitmq.GenerateUploadMQParam(music_file.UserID, music_file.UUID, music_file.FilePath)
	rabbitmq.RMQUpload.Publish(data)

	return music_file, true
}
