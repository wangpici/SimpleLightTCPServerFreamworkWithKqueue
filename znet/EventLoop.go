package znet

import (
	"fmt"
	"syscall"
	"test/iface"
)

type EventLoop struct {
	KqueueFd           int
	ListenerSocketUnit *SocketUnit
	SUManager          *SocketUnitManager
	MsgHandler         iface.IMessageHandler

	OnSocketUnitStart func(su iface.ISocketUnit)

	OnSocketUnitStop func(su iface.ISocketUnit)
}

func NewEventLoop(s *SocketUnit, server *Server) (*EventLoop, error) {
	kqueue, err := syscall.Kqueue()
	if err != nil {
		fmt.Println("kqueue init error, ", err)
		return nil, err
	}
	changeEvent := syscall.Kevent_t{
		Ident:  uint64(s.Fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	changeEventRegister, err := syscall.Kevent(kqueue, []syscall.Kevent_t{changeEvent}, nil, nil)
	if err != nil || changeEventRegister == 1 {
		fmt.Println(" register listen event error, ", err)
		return nil, err
	}

	return &EventLoop{
		KqueueFd:           kqueue,
		ListenerSocketUnit: s,
		SUManager:          NewSocketUnitManager(),
		MsgHandler:         server.MsgHandler,
	}, nil

}

func (e *EventLoop) Looping() {
	for {
		newEvent := make([]syscall.Kevent_t, syscall.SOMAXCONN)
		newEventCount, err := syscall.Kevent(e.KqueueFd, nil, newEvent, nil)
		if err != nil {
			fmt.Println(" loop error, ", err)
			continue
		}
		for i := 0; i < newEventCount; i++ {
			curEvent := newEvent[i]
			curEventFd := int(curEvent.Ident)

			if curEvent.Flags&syscall.EV_EOF != 0 {
				curSU, _ := e.SUManager.Get(curEventFd)
				curSU.Stop()

			} else if curEventFd == e.ListenerSocketUnit.Fd {
				newComingSocket, sa, err := syscall.Accept(curEventFd)
				if err != nil {
					fmt.Println(" newComingSocket accept error, ", err)
					continue
				}
				fmt.Println(" accept succ ")

				newComingEvent := syscall.Kevent_t{
					Ident:  uint64(newComingSocket),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
					Fflags: 0,
					Data:   0,
					Udata:  nil,
				}
				newComingEventReg, err := syscall.Kevent(e.KqueueFd, []syscall.Kevent_t{newComingEvent}, nil, nil)
				if err != nil || newComingEventReg == -1 {
					fmt.Println("reg newComingEvent error, ", err)
					continue
				}

				newComingSu := &SocketUnit{
					Fd:         newComingSocket,
					IsOpen:     true,
					Addr:       sa,
					ReadBuffer: NewByteArray(1024),
					MsgHandler: e.MsgHandler,
					SUManager:  e.SUManager,
					exitChan:   make(chan bool),
					msgChan:    make(chan []byte),
					ELoop:      e,
				}
				e.SUManager.Add(newComingSocket, newComingSu)
				newComingSu.Start()

			} else if curEvent.Filter&syscall.EVFILT_READ != 0 {
				//TODO 读取数据

				curSocketUnit, err := e.SUManager.Get(curEventFd)
				if err != nil {
					fmt.Println("cureventfd: ", curEventFd, " get su error, ", err)
				}
				fmt.Println(" get curSU succ, sufd: ", curEventFd)
				fmt.Println(" start process stream ")

				curSocketUnit.OnDataReceive()

			}

		}
	}
}

func (e *EventLoop) Stop() {
	e.SUManager.Clear()
}

func (e *EventLoop) SetOnSocketUnitStart(startHookFunc func(su iface.ISocketUnit)) {
	e.OnSocketUnitStart = startHookFunc
}

func (e *EventLoop) SetOnSocketUnitStop(stopHookFunc func(su iface.ISocketUnit)) {
	e.OnSocketUnitStop = stopHookFunc
}

func (e *EventLoop) CallOnSocketUnitStart(su iface.ISocketUnit) {
	if e.OnSocketUnitStart == nil {
		return
	}
	e.OnSocketUnitStart(su)
}

func (e *EventLoop) CallOnSocketUnitStop(su iface.ISocketUnit) {
	if e.OnSocketUnitStop == nil {
		return
	}
	e.OnSocketUnitStop(su)
}

func (e *EventLoop) StopAllSocketUnit() {
	e.SUManager.Clear()
}
