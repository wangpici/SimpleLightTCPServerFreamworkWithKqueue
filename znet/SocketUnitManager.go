package znet

import (
	"errors"
	"test/iface"
)

type SocketUnitManager struct {
	SocketUnits map[int]iface.ISocketUnit
}

func NewSocketUnitManager() *SocketUnitManager {
	return &SocketUnitManager{
		SocketUnits: make(map[int]iface.ISocketUnit),
	}
}

func (sum *SocketUnitManager) Add(fd int, su iface.ISocketUnit) error {
	if sum.SocketUnits[fd] != nil {
		return errors.New(" socketfd already exist! ")
	}

	sum.SocketUnits[fd] = su
	return nil
}

func (sum *SocketUnitManager) Remove(fd int) error {
	if sum.SocketUnits[fd] == nil {
		return errors.New(" socketfd not exist ")
	}
	delete(sum.SocketUnits, fd)
	return nil
}

func (sum *SocketUnitManager) Get(fd int) (iface.ISocketUnit, error) {
	if sum.SocketUnits[fd] == nil {
		return nil, errors.New(" get socketfd error, fd not exist! ")
	}
	return sum.SocketUnits[fd], nil
}

func (sum *SocketUnitManager) Clear() {
	for key, value := range sum.SocketUnits {
		delete(sum.SocketUnits, key)
		value.Stop()
	}
}
