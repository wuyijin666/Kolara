package main

import "Kolara/knet"

// 利用Kolara框架的api, 开发服务器端应用程序
func runServer() {
	// 1. 创建一个server句柄，利用Kolara框架的api
	s := knet.NewServer("Kolara V0.2")

	// 2. 启动服务器
	s.Serve()
}

