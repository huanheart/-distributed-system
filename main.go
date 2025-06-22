package main

import (
	"MyChat/common/mysql"
	"MyChat/common/redis"
	"MyChat/config"
	"MyChat/router"
	"fmt"
	"log"
)

func StartServer(addr string, port int) error {
	r := router.InitRouter()
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

func main() {
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	//初始化mysql
	if err := mysql.InitMysql(); err != nil {
		log.Println("InitMysql error , " + err.Error())
		return
	}
	//初始化redis
	redis.Init()

	err := StartServer(host, port) // 启动 HTTP 服务，监听 8080 端口
	if err != nil {
		panic(err)
	}

}
