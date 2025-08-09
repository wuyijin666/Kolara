package kiface

/**
 *  消息管理抽象层
**/

type IMsgHandle interface{
	DoMsgHandle(request IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}