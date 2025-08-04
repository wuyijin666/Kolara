package main

import(
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client is starting...")
	// 连接远程服务器
	conn, err := net.Dial("tcp" , "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start is err :" , err)
		return
	}

	for {
		// 连接调用write 写数据
		_, err := conn.Write([]byte("hello im KolaraV0.2 !"))
		if err != nil {
			fmt.Println("write to server err: ", err)
			return
		}

		buf := make([]byte, 512)
	    cnt, err :=	conn.Read(buf)
		if err != nil {
			fmt.Println("conn read is err :" , err)
			return 
		}
		fmt.Printf("server call back : %s, cnt = %d\n", buf[:cnt], cnt)

		// cpu阻塞， 目的是不要一直for循环
		time.Sleep(1*time.Second)

	}
}