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

func (this *PingRouter)	PreHandle(request kiface.IRequest) {
	fmt.Println("Call PreHandle Router")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("call before ping err")
	}
}

func (this *PingRouter)	Handle(request kiface.IRequest) {
	fmt.Println("Call Handle Router")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call ping err")
	}
}

func (this *PingRouter) PostHandle(request kiface.IRequest){
	fmt.Println("Call PostHandle Router")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping\n"))
	if err != nil {
		fmt.Println("call post ping err")
	}
}


func runServer() {
	// 1. 创建一个server句柄，利用Kolara框架的api
	s := knet.NewServer("Kolara V0.3")

	// 2. 给当前框架 添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 3. 启动服务器
	s.Serve()
}


func main() {
	runServer()
}
