package ziface

type IConnManager interface {
	Add(IConnection)
	Remove(IConnection)
	Get(uint32) (IConnection, error)
	Len() int
	CleanConn()
}
