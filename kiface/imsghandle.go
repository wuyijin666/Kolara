package kiface

type IMsgHandle interface{
	DoMsgHandle(request IRequest)
	AddRouter(msgId uint32, router IRouter)

}