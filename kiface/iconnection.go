package kiface

import "net"

// 定义连接模块抽象层
type IConnection interface {
	Start()

	Stop()

	// 获取当前连接绑定的conn (套接字)  
	GetTCPConnection() *net.TCPConn 

	// 获取当前连接模块的连接id
	GetConnID() uint32

	// 获取当前连接模块的客户端的TCP状态 ip port
	RemoteAddr() net.Addr

	// 发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error

}

// 连接绑定的处理函数的业务类型
type HandleFunc func(*net.TCPConn, []byte, int) error