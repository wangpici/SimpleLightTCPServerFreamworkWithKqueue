package iface

type ISocketUnitManager interface {
	Add(fd int, su ISocketUnit) error

	Remove(fd int) error

	Get(fd int) (ISocketUnit, error)

	Clear()
}
