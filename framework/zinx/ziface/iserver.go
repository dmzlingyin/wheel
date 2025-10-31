package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
}
