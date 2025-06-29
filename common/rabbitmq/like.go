package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type LikeMQParam struct {
	UserID    int64  `json:"user_id"`
	Action    int64  `json:"action"`
	LikeCount int64  `json:"like_count"`
	MusicUUID string `json:"uuid"`
}

// GenerateLikeMQParam 生成传入 Like MQ 的参数
func GenerateLikeMQParam(user_id int64, Action int64, LikeCount int64, MusicUUID string) []byte {
	param := LikeMQParam{
		UserID:    user_id,
		Action:    Action,
		LikeCount: LikeCount,
		MusicUUID: MusicUUID,
	}
	data, _ := json.Marshal(param)
	return data
}

func UpdateLikeCount(msg *amqp.Delivery) error {
	var param LikeMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	//更新music_file的表
	mysql.UpdateLikeCount(param.LikeCount, param.MusicUUID)
	return nil
}

// 这里还需要一个更新reaction表的表
func UpdateAction(msg *amqp.Delivery) error {
	var param LikeMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	//todo:这里需要先查一下是否有这个表，决定插入还是更新操作
	//todo:或者直接将插入和更新分开(我更倾向这种，因为外部已经判断了）
	mysql.UpdateAction(param.Action, param.UserID, param.MusicUUID)
	return nil
}

func InsertAction(msg *amqp.Delivery) error {
	var param LikeMQParam
	err := json.Unmarshal(msg.Body, &param)
	if err != nil {
		return err
	}
	//todo:这里需要先查一下是否有这个表，决定插入还是更新操作
	//todo:或者直接将插入和更新分开(我更倾向这种，因为外部已经判断了）
	mysql.InsertAction(param.Action, param.UserID, param.MusicUUID)
	return nil
}
