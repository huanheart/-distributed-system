package rabbitmq

var (
	RMQUpload        *RabbitMQ
	RMQRemoteUpload  *RabbitMQ
	RMQCountDuration *RabbitMQ
)

func InitRabbitMQ() {
	// 创建MQ并启动消费者
	// 无论调用多少次 NewWorkRabbitMQ，只会创建一次连接
	// 不同队列共用一个连接，可以保持不同队列消费消息的顺序
	RMQUpload = NewWorkRabbitMQ("Upload")
	go RMQUpload.Consume(Upload)

	RMQRemoteUpload = NewWorkRabbitMQ("RemoteUpload")
	go RMQRemoteUpload.Consume(RemoteUpload)

	RMQCountDuration = NewWorkRabbitMQ("CountDuration")
	go RMQCountDuration.Consume(CountDuration)
}

// DestroyRabbitMQ 销毁RabbitMQ

func DestroyRabbitMQ() {
	RMQUpload.Destroy()
	RMQRemoteUpload.Destroy()
}
