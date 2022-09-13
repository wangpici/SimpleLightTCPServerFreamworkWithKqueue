package iface

type IMessage interface {
	GetMsgId() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgId(id uint32)
	SetMsgData(data []byte)
	SetDataLen(len uint32)
}
