package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"wheel/framework/zinx/utils"
	"wheel/framework/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	buf := bytes.NewReader(binaryData)
	msg := &Message{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg")
	}
	return msg, nil
}
