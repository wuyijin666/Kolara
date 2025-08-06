package kiface

// request的抽象类：实际是把 客户端的请求链接 和 请求数据 包装到了一个request里

type IRequest interface {
	// 链接
	GetConnection() IConnection

	// 得到请求数据
	GetData() []byte

	GetMsgId() uint32
}