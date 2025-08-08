package knet

import (
  "Kolara/kiface"
  "fmt"
  "net"
  "Kolara/utils"
  
)


type Server struct { 
   Name string
   IP string
   IPVer string
   Port int

//    // 给当前的Server添加一个router,server注册的连接对应的处理业务
//    Router kiface.IRouter
   MsgHandle kiface.IMsgHandle
}

// // 定义当前连接所绑定的handle api (目前该handle写死，之后用户可自定义handle)
// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//      // 回显业务
// 	 fmt.Println("[Conn handle] CallbackToClient is called ...")
// 	 if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back buf err: ", err)
// 		return errors.New("CallbackToClient err")
// 	}
// 	return nil
// }

func (s *Server) Start() { 
	fmt.Printf("[Kolara] server name : %s is running at port %d ip %s\n", utils.GlobalObject.Name, utils.GlobalObject.TcpPort, utils.GlobalObject.Host)
    fmt.Printf("[Kolara] server version : %s\n", utils.GlobalObject.Version)


	go func(){

	// 1. 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVer, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	// 2. 监听服务器地址
	listener, err := net.ListenTCP(s.IPVer, addr)
	if err != nil {
		fmt.Println("listen tcp err: ", err)
		return
	}

	var cid uint32 = 0

	// 3. 阻塞等待客户端连接，处理客户端连接业务(读写)
	for {
		// 如果客户端连接，阻塞会返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept tcp err: ", err)
			continue
		}

		// 将处理新连接的业务方法与 conn 进行绑定，得到我们的连接模块
		dealConn := NewConnection(conn, cid, s.MsgHandle)
		cid ++

		// 启动当前连接的业务处理
		go dealConn.Start()
    }
}()
}

func (s *Server) Stop() { 

}

func (s *Server) Serve() { 
	s.Start()

	select{}

}

// 路由功能：给当前连接添加路由方法，供客户端连接处理使用
func(s *Server) AddRouter(msgId uint32, router kiface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("add router succ")
}
 

func NewServer(name string) kiface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IP: utils.GlobalObject.Host,
		IPVer: "tcp4",
		Port: utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
	}
	return s
}
