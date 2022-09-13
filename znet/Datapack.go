package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"test/iface"

	"test/utils"
)

type Datapack struct {
}

func (dp *Datapack) GetHeadLen() uint32 {

	// len(uint32) 4 bytes, id(uint32) 4 bytes
	return 8
}

func (dp *Datapack) Pack(msg iface.IMessage) ([]byte, error) {

	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *Datapack) Unpack(data []byte) (iface.IMessage, error) {

	dataBuff := bytes.NewReader(data)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too long msg recv")
	}
	return msg, nil
}
