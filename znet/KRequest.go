package znet

import "test/iface"

type KRequest struct {
	Su  iface.ISocketUnit
	Msg iface.IMessage
}

func (kr *KRequest) GetSocketUnit() iface.ISocketUnit {
	return kr.Su
}

func (kr *KRequest) GetData() []byte {
	return kr.Msg.GetData()
}

func (kr *KRequest) GetMsgId() uint32 {
	return kr.Msg.GetMsgId()
}
