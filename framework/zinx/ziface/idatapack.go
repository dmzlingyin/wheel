package ziface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(message IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
