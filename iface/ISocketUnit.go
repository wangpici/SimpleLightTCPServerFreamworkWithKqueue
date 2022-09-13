package iface

import "syscall"

type ISocketUnit interface {
	OnDataReceive()
	GetSocketFd() int
	SendMsg(msg IMessage) error
	StartWriteGoroutine()
	Start()
	Stop()
	OpenStatus() bool
	GetSockAddr() syscall.Sockaddr
	GetSUManager() ISocketUnitManager
	GetEventLoop() IEventLoop
}
