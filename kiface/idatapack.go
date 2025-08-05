package kiface

/**
* IDatapack 提供了解决TCP粘包问题的封包与拆包接口
* header包含msgLen msgId 
*/

type IDatapack interface {
    // Pack 将消息体长度、消息ID和数据打包成二进制流
    Pack(message IMessage) ([]byte, error)
    // Unpack 解析数据包，返回消息体长度和错误信息
    Unpack([]byte) (IMessage, error)
    // GetHeaderLen 返回协议头部的固定长度
    GetHeaderLen() uint32
}