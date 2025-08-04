package knet

import (
	"Kolara/kiface"
	"fmt"
	"net"
)

// 定义连接属性
type Connection struct {
	// 当前连接的TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态
    isClosed bool
	// 当前连接所绑定的处理业务api
	handleAPI kiface.HandleFunc

	// 告知当前连接已经停止/退出的channel
	ExitChan chan bool

	// 该连接处理的方法Router
	Router kiface.IRouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, router kiface.IRouter) *Connection {
	c := &Connection {
		Conn : conn,
		ConnID : connId,
		isClosed : false,
		Router: router,
		ExitChan : make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is starting ... ")
	defer fmt.Println("Conn is closing ... ConnId = ", c.ConnID, "RemoteAddr = " , c.RemoteAddr().String())
	defer c.Stop()

	for{

		buf := make([]byte, 4096)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read from conn err: ", err)
			continue
		}

		// // 调用当前连接绑定的handleAPI (后续被框架集成Router模块取代)
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("connID: " , c.ConnID, "handle is err: ", err)
		// 	break;
		// }

		// 得到当前conn的Request请求数据
		req := Request {
			conn : c,
			data : buf,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		// 执行注册的路由方法
		go func(request kiface.IRequest){
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)

		}(&req)
		
	}
}


// 启动连接 让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("conn Start() ... ConnId = " , c.ConnID)
	// 连接进行读业务
	go c.StartReader()

	// 连接进行写业务
	

	


}


func (c *Connection) Stop() {

}

// 获取当前连接绑定的Socket 
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn

}

// 获取当前连接模块的连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID

}

// 获取当前连接模块的客户端的TCP状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil

}
