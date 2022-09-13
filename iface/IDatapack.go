package iface

type IDatapack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(data []byte) (IMessage, error)
}
