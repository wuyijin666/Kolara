package main

import (
	"Kolara/knet"
	"Kolara/kiface"
	"fmt"
)

// 利用Kolara框架的api, 开发服务器端应用程序

// ping test 自定义路由
type PingRouter struct {
	knet.BaseRouter
}



func (this *PingRouter)	Handle(request kiface.IRequest) {
	fmt.Println("Call Handle Router")
	// 先读取客户端数据 再回写ping ping ping
	fmt.Println("recv from client msgId:", request.GetMsgId(), "data:", string(request.GetData()))

    err := request.GetConnection().SendMsg(1, []byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call ping err")
	}
}



func runServer() {
	// 1. 创建一个server句柄，利用Kolara框架的api
	s := knet.NewServer("Kolara V0.5")

	// 2. 给当前框架 添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3. 启动服务器
	s.Serve()
}


// func main() {
// 	runServer()
// }
