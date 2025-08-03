package knet

import (
	"Kolara/kiface"
)

type Request struct {
	conn kiface.IConnection
	data []byte
}

func (r *Request) GetConnection() kiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}