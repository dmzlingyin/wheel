package ziface

type IRouter interface {
	PreHandle(IRequest)
	Handle(IRequest)
	PostHandle(IRequest)
}
