package knet

import (
	"Kolara/kiface"
	"fmt"
	"net"
	"io"
	"errors"
	"Kolara/utils"
	"sync"
)

// 定义连接属性
type Connection struct {
	TcpServer kiface.IServer
	// 当前连接的TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态
    isClosed bool
	// 当前连接所绑定的处理业务api
	handleAPI kiface.HandleFunc

	// 告知当前连接已经停止/退出的channel (由Reader告知Writer该退出信号)
	ExitChan chan bool

	// 无缓冲管道， 用于读写Goroutine之间的消息通信
	MsgChan chan []byte

	// 该连接处理的方法Router
	// Router kiface.IRouter

	// 消息管理：msgId和对应的处理业务api之间的关系
	MsgHandle kiface.IMsgHandle

	// 连接属性集合
	Property map[string]interface{}
    // 保护连接属性的锁
	PropertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server kiface.IServer, conn *net.TCPConn, connId uint32, msgHandle kiface.IMsgHandle) *Connection {
	c := &Connection {
		TcpServer: server,
		Conn : conn,
		ConnID : connId,
		isClosed : false,
		ExitChan : make(chan bool, 1),
		MsgChan : make(chan []byte),
		MsgHandle: msgHandle,
		Property: make(map[string]interface{}),
	}

	// 将conn加入到ConnMgr中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is starting]")
	defer fmt.Println("Conn Reader close, ConnId = ", c.ConnID, "RemoteAddr = " , c.RemoteAddr().String())
	defer c.Stop()

	for{
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("read from conn err: ", err)
		// 	continue
		// }
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的Msg Head, 二进制流， 8字节
		headData := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read head error", err)
			break
		}
		// 拆包 得到msgId 和 msgDataLen, 放到msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据msgDataLen 读取Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read data error", err)
				break
			}
		}
		msg.SetData(data)

		// // 调用当前连接绑定的handleAPI (后续被框架集成Router模块取代)
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("connID: " , c.ConnID, "handle is err: ", err)
		// 	break;
		// }

		// 得到当前conn的Request请求数据
		req := Request {
			conn : c,
			msg : msg,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		// 执行注册的路由方法
		// go func(request kiface.IRequest){
		// 	c.Router.PreHandle(request)
		// 	c.Router.Handle(request)
		// 	c.Router.PostHandle(request)

		// }(&req)
		
		if utils.GlobalObject.WorkPoolSize > 0 {
			// 已开启worker工作池
			c.MsgHandle.SendMsgToTaskQueue(&req)
		}else {
			// 根据msgId调用对应的router业务
		    go c.MsgHandle.DoMsgHandle(&req)
		}	
	}
}

// 写Goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is starting]")
	defer fmt.Println(c.RemoteAddr().String(), "conn Writer exit now.")

	// 不断的阻塞等待channel的消息， 进行写给客户端
	for {
		select {
		case data := <- c.MsgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:, ", err)
				return 
			}
		case <- c.ExitChan:
			// 代表Reader已退出，则Writer也退出。conn已经关闭
			return
		}
	}
}


// 启动连接 让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("conn Start() ... ConnId = " , c.ConnID)
	// 连接进行读业务
	go c.StartReader()
	// 启动当前连接进行写业务
    go c.StartWriter()	

	c.TcpServer.CallOnConnStart(c)
}


func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnId = ", c.ConnID)
	// 若当前连接已关闭
    if c.isClosed == true {
		return
	} 
	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 释放连接资源
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)

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
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("conn is closed")
	}
	// 封包 msgId|msgLen|Data
	dp := NewDataPack() 
	binaryData, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		return errors.New("pack err ...")
	}
	// 将数据发送给客户端
	// _, err = c.Conn.Write(binaryData)
	// if err != nil {
	// 	fmt.Println("write msgId", msgId, "err",  err)
	// 	return errors.New("conn Write err...")
	// }

	// 将封包的消息发送给管道
	c.MsgChan <- binaryData
	return nil
}

func(c *Connection) AddProperty(key string, value interface{}){
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	
	c.Property[key] = value
}

func(c *Connection) GetProperty(key string) (interface{}, error) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	if value, ok := c.Property[key]; ok {
		return value, nil
	}
	return nil, errors.New("property not found")	
}

func(c *Connection) RemoveProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	delete(c.Property, key)
}
