package rabbitmq

import (
	"MyChat/common/mysql"
	"encoding/json"
	"github.com/streadway/amqp"
)

type UploadMQParam struct {
	UserID   int64  `json:"user_id"`
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

// GenerateLikeMQParam 生成传入 Like MQ 的参数
func GenerateUploadMQParam(user_id int64, file_id, file_path string) []byte {
	param := UploadMQParam{
		UserID:   user_id,
		FileID:   file_id,
		FilePath: file_path,
	}
	data, _ := json.Marshal(param)
	return data
}

// 将本地的false->true
// 直接将mysql中file_path对应的行所对应的字段的is_upload置为true
// 用不到file_id和user_id
func Upload(msg *amqp.Delivery) error {
	var param UploadMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	mysql.MarkMusicFileUploaded(param.FilePath, 1)
	return nil
}
