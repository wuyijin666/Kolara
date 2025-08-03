package knet
import (
	"Kolara/kiface"
)

//实现router, 先嵌入BaseRouter这个基类，之后按需对BaseRouter进行重写

type BaseRouter struct {}

func (br *BaseRouter) PreHandle(request kiface.IRequest) {}

func (br *BaseRouter) Handle(request kiface.IRequest) {}

func (br *BaseRouter) PostHandle(request kiface.IRequest) {}