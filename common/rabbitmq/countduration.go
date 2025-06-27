package rabbitmq

import (
	"MyChat/common/ffmpeg"
	"MyChat/common/mysql"
	"encoding/json"
	"github.com/streadway/amqp"
)

type CDMQParam struct {
	FilePath string `json:"file_path"`
}

func GenerateCDMQParam(file_path string) []byte {
	param := CDMQParam{
		FilePath: file_path,
	}
	data, _ := json.Marshal(param)
	return data
}

// 将本地的false->true
// 直接将mysql中file_path对应的行所对应的字段的is_upload置为true
// 用不到file_id和user_id
func CountDuration(msg *amqp.Delivery) error {
	var param CDMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	val, err := ffmpeg.CountDuration(param.FilePath)
	if err != nil {
		return err
	}
	mysql.SetCountDuration(param.FilePath, val)
	return nil
}
