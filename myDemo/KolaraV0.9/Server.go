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

type HelloKolaraRouter struct {
	knet.BaseRouter
}

func (p *PingRouter) Handle(request kiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	// 先读取客户端数据 再回写ping ping ping
	fmt.Println("recv from client msgId:", request.GetMsgId(), "data:", string(request.GetData()))

    err := request.GetConnection().SendMsg(200, []byte("ping ping ping\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HelloKolaraRouter) Handle(request kiface.IRequest) {
	fmt.Println("Call HelloKolaraRouter Handle")
	// 先读取客户端数据 再回写ping ping ping
	fmt.Println("recv from client msgId:", request.GetMsgId(), "data:", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello welcome to KolaraV0.8!\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnBeginHook(conn kiface.IConnection){
	fmt.Println("DoConnBeginHook is called ...")
	if err := conn.SendMsg(202, []byte("DOCONN BEGIN...")); err != nil {
		fmt.Println(err)
	}

}

func DoConnStopHook(conn kiface.IConnection){
	fmt.Println("DoConnStopHook is called ...")
}



func runServer() {
	// 1. 创建一个server句柄，利用Kolara框架的api
	s := knet.NewServer("Kolara V0.9")

	// 注册hook函数 实现用户方的一些自定义逻辑
	s.SetOnConnStart(DoConnBeginHook)
	s.SetOnConnStop(DoConnStopHook)

	// 2. 给当前框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloKolaraRouter{})

	// 3. 启动服务器
	s.Serve()
}


// func main() {
// 	runServer()
// }
