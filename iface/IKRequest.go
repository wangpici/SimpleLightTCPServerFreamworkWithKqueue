package iface

type IKRequest interface {
	GetSocketUnit() ISocketUnit
	GetData() []byte
	GetMsgId() uint32
}
