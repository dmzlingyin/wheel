package znet

import "wheel/framework/zinx/ziface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (br *BaseRouter) Handle(request ziface.IRequest) {
}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {
}
