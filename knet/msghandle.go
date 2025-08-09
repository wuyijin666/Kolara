package knet

import (
	"Kolara/kiface"
	"Kolara/utils"
	"fmt"
)

type MsgHandle struct {
	// 存放每个msgId对应的处理方法
	Apis map[uint32] kiface.IRouter
	// 负责从Worker中取任务的消息队列
	TaskQueue []chan kiface.IRequest
	// 业务Worker池的worker数量
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] kiface.IRouter),	
		TaskQueue: make([]chan kiface.IRequest, utils.GlobalObject.WorkPoolSize),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize,
	}
}

// 针对msgId 执行/调度特定的router处理方法
func (mh *MsgHandle) DoMsgHandle(request kiface.IRequest) {
	router, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("msgId: %d not registered\n", request.GetMsgId())
		return
	}
	
	router.PreHandle(request)
    router.Handle(request)
	router.PostHandle(request)
}

// 为消息msgId添加具体的路由处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router kiface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		fmt.Printf("repeated api, msgId = %d\n ",  msgId)
		return 
    }
	mh.Apis[msgId] = router
}

// 启动一个Worker工作池 (一个框架只能有一个worker池)
func (mh *MsgHandle) StartWorkerPool() {	
	// 根据WorkPoolSize去分别开启worker, 每个worker用一个go承载
	for i := 0; i< int(mh.WorkPoolSize); i++ {
		// 一个worker被启动
		// 1. 当前的worker对应的消息队列channel， 第0个worker 就对应 第0个channel
		mh.TaskQueue[i] = make(chan kiface.IRequest, utils.GlobalObject.WorkPoolSize)
		// 2. 启动当前worker, 阻塞等待消息从channel传递进来
		go mh.StartOneWorker(uint32(i), mh.TaskQueue[i])
	}

}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workId uint32, taskQueue chan kiface.IRequest) {
	fmt.Printf("Worker ID = %d is started\n", workId)
	for {
		// 不停阻塞等待对应消息队列的消息
		select {
		case request := <- taskQueue:
			// 若有消息过来，出列的就是一个客户端的request, 执行当前request所绑定的业务
			mh.DoMsgHandle(request)

		}
	}
}

// 将消息交给TaskQueue,由对应的worker处理 
func (mh *MsgHandle) SendMsgToTaskQueue(request kiface.IRequest) {
	// 将消息平均分配给不同的worker 
	// 根绝客户端建立的connId来进行分配
	workId := request.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println("Add connId", request.GetConnection().GetConnID(), 
                "request msgId", request.GetMsgId(),
			    "to workId" , workId)
	mh.TaskQueue[workId] <- request
}