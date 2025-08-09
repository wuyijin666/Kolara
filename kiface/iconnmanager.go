package kiface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	// 根据连接id， 获取连接
	Get(connId uint32) (IConnection, error)
	// 获取当前连接总数
	Len() int
	// 清理所有连接
	ClearConn()
}