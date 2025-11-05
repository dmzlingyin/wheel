package znet

import (
	"errors"
	"sync"
	"wheel/framework/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	lock        sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.connections[conn.GetConnID()] = conn
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.connections, conn.GetConnID())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) CleanConn() {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
}
