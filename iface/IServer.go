package iface

type IServer interface {
	//Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, router IRouter)
	StartKqueueServer()
	//	GetConnectionManager() IConnectionManager
	SetOnSocketUnitStart(startHookFunc func(su ISocketUnit))
	SetOnSocketUnitStop(stopHookFunc func(su ISocketUnit))
	SetOnEventLoopStart(startHookFunc func(el IEventLoop))
	SetOnEventLoopStop(stopHookFunc func(el IEventLoop))

	CallOnEventLoopStart(el IEventLoop)

	CallOnEventLoopStop(el IEventLoop)
}
