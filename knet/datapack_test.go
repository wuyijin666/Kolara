package knet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	/** 
	* 模拟服务端
	*/ 
	// 1. 与客户端建立连接
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err ...", err)
		return 
	}
	// 2. 监听listener, 获取conn
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error...", err)
				continue
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// --- 拆包 ---
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
						// 2.1 从conn中拆包 获取header
				headData := make([]byte, dp.GetHeaderLen())
				if _, err := io.ReadFull(conn, headData); err != nil {
					fmt.Println("read header data err:" , err)
					return 
				}

				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return 
				}
				// 2.2 拆包，根据msgLen,读取data消息体
				if msgHead.GetMsgLen() > 0 {
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())

					// 读取dataLen的长度，再次从io流中获取
					if _, err := io.ReadFull(conn, msg.Data); err != nil {
						fmt.Println("data read false...", err) 
						return 
					}
  
					fmt.Printf("recv msgId: %d, msgLen: %d, data: %s\n", msg.GetMsgId(), msg.GetMsgLen(), string(msg.GetData()))
				}
			}
		}(conn)
	}
   }()
	
	



	/**
	*    模拟客户端
	**/
	// 1. 建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
	}
	// 2. 创建一个封包对象dp
	dp := NewDataPack()

	// 模拟粘包过程 封装两个msg一同发送
	msg1 := &Message{
		MsgId: 1,
		MsgLen: 5,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err:", err)
		return 
	}

	msg2 := &Message{
		MsgId: 2,
		MsgLen: 7,
		Data: []byte{'K', 'o', 'l', 'a', 'r','a','!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err:", err)
		return 
	}
 
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 将处理完的粘包 一起发送给服务端 
	conn.Write(sendData1)

    // // 客户端阻塞
	// select {}
	time.Sleep(50*time.Second)
}