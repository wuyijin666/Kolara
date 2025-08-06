package knet

import (
	"Kolara/kiface"
)

type Request struct {
	conn kiface.IConnection
	msg kiface.IMessage
}

func (r *Request) GetConnection() kiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}