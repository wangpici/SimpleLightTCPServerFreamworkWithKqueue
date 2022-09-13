package znet

import (
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
	"test/iface"
)

type SocketUnit struct {
	Fd         int
	IsOpen     bool
	Addr       syscall.Sockaddr
	ReadBuffer ByteArray
	MsgHandler iface.IMessageHandler
	msgChan    chan []byte
	exitChan   chan bool
	SUManager  *SocketUnitManager
	ELoop      iface.IEventLoop
}

func (su *SocketUnit) GetSocketFd() int {
	return su.Fd
}

func (su *SocketUnit) OpenStatus() bool {
	return su.IsOpen
}

func (su *SocketUnit) GetSockAddr() syscall.Sockaddr {
	return su.Addr
}

func (su *SocketUnit) GetSUManager() iface.ISocketUnitManager {
	return su.SUManager
}

func (su *SocketUnit) GetEventLoop() iface.IEventLoop {
	return su.ELoop
}

func Listen(ip string, port int) (*SocketUnit, error) {
	s := &SocketUnit{}
	socketFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("get socketfd error, ", err)
		return nil, err
	}
	socketAddr := &syscall.SockaddrInet4{
		Port: port,
	}
	copy(socketAddr.Addr[:], net.ParseIP(ip))

	if err = syscall.Bind(socketFd, socketAddr); err != nil {
		fmt.Println(" bind addr error, ", err)
		return nil, err
	}
	if err = syscall.Listen(socketFd, syscall.SOMAXCONN); err != nil {
		fmt.Println(" listen start error, ", err)
		return nil, err
	}

	s.Fd = socketFd
	s.IsOpen = true
	return s, nil

}

func (su *SocketUnit) OnDataReceive() {

	temp := make([]byte, 1024)
	n, _, err := syscall.Recvfrom(su.Fd, temp, syscall.MSG_DONTWAIT)
	recvBytes := temp[:n]
	if err != nil {
		fmt.Println(" recv error, socketunit: ", su.Fd, " error: ", err)
		return
	}
	fmt.Println(" recv ", string(recvBytes), " start divide and merge... ")

	su.ReadBuffer.Write(recvBytes)
	readIndex := su.ReadBuffer.ReadIndex()
	writeIndex := su.ReadBuffer.WriteIndex()
	count := 0
	buffCopy := su.ReadBuffer.Bytes()
	for {
		if su.ReadBuffer.BufferLen() <= 8 {
			break
		}

		lenBytes := buffCopy[readIndex : readIndex+4]
		bodyLen := binary.LittleEndian.Uint32(lenBytes)

		if int(bodyLen)+8 > writeIndex-readIndex {
			break
		}
		fmt.Println(" lenbytes: ", int(binary.LittleEndian.Uint32(lenBytes)))
		totalRecvLen := int(bodyLen) + 8
		count += totalRecvLen
		//msgBytes := su.ReadBuffer.Read(int(bodyLen) + 8)
		msgBytes := buffCopy[readIndex : readIndex+totalRecvLen]
		readIndex += totalRecvLen

		msgIdBytes := msgBytes[4:8]
		msgId := binary.LittleEndian.Uint32(msgIdBytes)
		msgData := msgBytes[8:]

		fmt.Println(" msgdata: ", string(msgData))
		fmt.Println(" readbuff cur len: ", writeIndex-readIndex)
		msg := NewMessage(msgId, msgData)
		krequest := &KRequest{
			Su:  su,
			Msg: msg,
		}
		su.MsgHandler.SendRequestToTaskQue(krequest)
	}
	su.ReadBuffer.ReadWithoutReturn(count)
}

func (su *SocketUnit) SendMsg(msg iface.IMessage) error {
	dp := Datapack{}
	sendBytes, err := dp.Pack(msg)
	if err != nil {
		fmt.Println(" pack msgid: ", string(msg.GetMsgId()), " error: ", err)
		return err
	}
	//TODO 塞进管道
	su.msgChan <- sendBytes
	return nil
}

func (su *SocketUnit) StartWriteGoroutine() {
	fmt.Println(" socketunit: ", string(su.Fd), " write goroutine start...")
	defer fmt.Println(" socketunit: ", string(su.Fd), " write goroutine exist...")

	for {
		select {
		case data := <-su.msgChan:
			err := syscall.Sendto(su.Fd, data, syscall.MSG_DONTWAIT, su.Addr)
			if err != nil {
				fmt.Println(" send msg: ", string(data), " error: ", err)
			}
		case <-su.exitChan:
			return
		}
	}

}

func (su *SocketUnit) Stop() {

	if su.IsOpen != false {
		return
	}

	su.IsOpen = false
	su.exitChan <- true
	su.SUManager.Remove(su.Fd)

	// TODO： 退出的钩子函数还要调用一下
	su.ELoop.CallOnSocketUnitStop(su)
	close(su.exitChan)
	close(su.msgChan)
	syscall.Close(su.Fd)
}

func (su *SocketUnit) Start() {
	fmt.Println(" socketunit: ", string(su.Fd), " start working ")

	go su.StartWriteGoroutine()
	//TODO: 开始的钩子函数调用下
	su.ELoop.CallOnSocketUnitStart(su)

}
