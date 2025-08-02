package knet

import (
  "Kolara/kiface"
  "fmt"
  "net"
)


type Server struct { 
   Name string
   IP string
   IPVer string
   Port int
}

func (s *Server) Start() { 

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

	// 3. 阻塞等待客户端连接，处理客户端连接业务(读写)
	for {
		// 如果客户端连接，阻塞会返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept tcp err: ", err)
			continue
		}

		// 4. 启动一个goroutine处理连接 先做基础的业务，最大允许512字节的回显操作
		go func() {
			// 读取客户端的数据
			for {
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Println("read from conn err: ", err)
					continue
				}

				// 回显数据
				if _, err := conn.Write(buf[:cnt]); err != nil {
					fmt.Println("write to conn err: ", err)
					continue
				}
			}
		}()
    }
}()
}

func (s *Server) Stop() { 

}

func (s *Server) Serve() { 
	s.Start()

	select{}

}

func NewServer(name string) kiface.Server {
	s := &Server{
		Name: name,
		IP: "0.0.0.0",
		IPVer: "tcp4",
		Port: 8999,
	}
	return s
}
