package iface

type IEventLoop interface {
	Stop()
	Looping()

	SetOnSocketUnitStart(startHookFunc func(su ISocketUnit))

	SetOnSocketUnitStop(stopHookFunc func(su ISocketUnit))

	CallOnSocketUnitStart(su ISocketUnit)

	CallOnSocketUnitStop(su ISocketUnit)
}
