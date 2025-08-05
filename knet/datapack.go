package knet

import (
	"bytes"
	"Kolara/kiface"
	"encoding/binary"
	"errors"
	"Kolara/utils"
)

// 封包 拆包的具体模块
type DataPack struct {}

func NewDataPack() *DataPack {
	return &DataPack{}
}


// 封包
func(dp *DataPack) Pack(message kiface.IMessage) ([]byte, error){
  // 创建一个存放bytes字节的缓冲
  dataBuf := bytes.NewBuffer([]byte{})
  // 将msgLen 读取到buf缓冲区
if err := binary.Write(dataBuf, binary.LittleEndian, message.GetMsgLen()); err != nil {
	return nil, err
}

  // 将msgId 读取到buf缓冲区
if err := binary.Write(dataBuf, binary.LittleEndian, message.GetMsgId()); err != nil {
	return nil, err
}

  // 将消息体data读取到buf缓冲区
if err :=  binary.Write(dataBuf, binary.LittleEndian, message.GetData()); err != nil {
	return nil, err
}

   return dataBuf.Bytes(), nil
}

// 拆包
func(dp *DataPack) Unpack(binaryData []byte) (kiface.IMessage, error){
   //创建一个输入二进制数据的ioReader
   dataBuf := bytes.NewReader(binaryData)

   // 解析header 拿到msgLen和msgId
   msg := &Message{}

   // 读取msgLen
   if err := binary.Read(dataBuf, binary.LittleEndian, &msg.MsgLen); err != nil {
	return nil, err
   }

   // 读取msgId
   if err := binary.Read(dataBuf, binary.LittleEndian, &msg.MsgId); err != nil {
	return nil, err
   }

   if utils.GlobalObject.MaxPackageSize > 0 && msg.MsgLen > utils.GlobalObject.MaxPackageSize {
   return nil, errors.New("too large package recv")
   }
 
   return msg, nil
}

func(dp *DataPack) GetHeaderLen() uint32 {
	// msgLen uint32 (4字节) + msgId uint32 (4字节) = 8字节
	return 8 
}
