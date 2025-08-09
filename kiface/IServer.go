package kiface

type IServer interface {
	Start()
	Stop()
	Serve()

	// 路由功能: 给当前服务注册一个路由业务方法， 供客户端连接处理使用
	AddRouter(msgId uint32, router IRouter)
	GetConnMgr() IConnManager 
}	
