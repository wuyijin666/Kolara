package knet

import (
	"Kolara/kiface"
	"fmt"
)

type MsgHandle struct {
	// 存放每个msgId对应的处理方法
	Apis map[uint32] kiface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] kiface.IRouter),	
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

