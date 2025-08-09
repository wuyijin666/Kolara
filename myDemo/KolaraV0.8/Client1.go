package main

import (
	"Kolara/knet"
	"fmt"
	"net"
	"time"
	"io"
)

func main()() {
	fmt.Println("client1 is starting...")
	// 连接远程服务器
	conn, err := net.Dial("tcp" , "127.0.0.1:7890")
	if err != nil {
		fmt.Println("client1 start is err :" , err)
		return
	}

	for {
		// 发送封包的msg消息 MsgId : 1
		dp := knet.NewDataPack()
		binaryData, err := dp.Pack(knet.NewMsgPackage(1, []byte("hello Kolara V0.8 client1 test send msg!")))
		if err != nil {
			fmt.Println("pack err :" , err)
			return
		}
		// 发送数据
		_, err = conn.Write(binaryData)
		if err != nil {
			fmt.Println("write err :" , err)
			return
		}

	    // 服务器给我们回一个message MsgId: 1  Hello welcome to KolaraV0.6!
       // 1. 拆服务器的包 先拿header
	   binaryHead := make([]byte, dp.GetHeaderLen())
	   if _, err := io.ReadFull(conn, binaryHead); err != nil {
		 fmt.Println("read err", err)
		 break;
	   }
	
	   msgHead, err := dp.Unpack(binaryHead)
	   if err != nil {
		fmt.Println("unpack err", err)
		break
	   }
	   if msgHead.GetMsgLen() > 0 { 
		 msg := msgHead.(*knet.Message)
		 msg.Data = make([]byte, msgHead.GetMsgLen())

		if _, err := io.ReadFull(conn, msg.Data); err != nil { 
			fmt.Println("read msg data err", err)
			return
		}
	
	fmt.Println("----->recv msg: ", msgHead.GetMsgId(), ", len: ", msgHead.GetMsgLen(), ", data: ", string(msg.Data))  
	} 
		// cpu阻塞， 目的是不要一直for循环
		time.Sleep(1*time.Second)

	}
}