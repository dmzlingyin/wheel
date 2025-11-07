package znet

import (
	"fmt"
	"wheel/framework/zinx/zconf"
	"wheel/framework/zinx/ziface"
)

type MsgHandle struct {
	Apis         map[uint32]ziface.IRouter
	WorkPoolSize uint32
	TaskQueue    []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: zconf.GlobalObject.WorkerPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, zconf.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgID]; ok {
		panic(fmt.Sprintf("repeated api, msgID = %d", msgID))
	}
	m.Apis[msgID] = router
	fmt.Println("add api success: ", msgID)
}

func (m *MsgHandle) StartWorkPool() {
	for i := 0; i < int(m.WorkPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, zconf.GlobalObject.MaxWorkerTaskLen)
		go func(id int) {
			fmt.Println("worker ", id, " is started...")
			for {
				select {
				case request := <-m.TaskQueue[id]:
					fmt.Println("worker ", id, " is working...")
					m.DoMsgHandler(request)
				}
			}
		}(i)
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	id := request.GetConnection().GetConnID() % m.WorkPoolSize
	m.TaskQueue[id] <- request
}
