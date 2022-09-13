package main

import (
	"fmt"
	"test/iface"
	"test/znet"
)

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (h *HelloZinxRouter) Handle(krequest iface.IKRequest) {

	recvByte := krequest.GetData()

	recvString := " recv " + string(recvByte)
	fmt.Println(recvString)
	su := krequest.GetSocketUnit()
	replyMsg := znet.NewMessage(0, []byte(recvString))
	if err := su.SendMsg(replyMsg); err != nil {
		fmt.Println(" recv completed but send error ")
	}

}

var helloZinxRouter HelloZinxRouter

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &helloZinxRouter)

	s.Serve()
}
