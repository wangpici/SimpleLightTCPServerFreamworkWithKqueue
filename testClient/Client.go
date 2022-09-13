package main

import "fmt"

import (
	"io"
	"net"
	"test/znet"
	"time"
)

func main() {

	fmt.Println(" client start ")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "192.168.1.105:8999")
	if err != nil {
		fmt.Println(" client start err, ", err)
		return
	}

	time.Sleep(2 * time.Second)

	for true {

		dp := znet.Datapack{}
		msgByte, err := dp.Pack(znet.NewMessage(0, []byte(" hello from client by tcp ")))
		if err != nil {
			fmt.Println(" pack msg error ", err)
			return
		}

		_, err = conn.Write(msgByte)
		if err != nil {
			fmt.Println(" write msg error ", err)
		}
		fmt.Println(" msg send succ ")

		byteHead := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(conn, byteHead); err != nil {
			fmt.Println(" read msg head error ", err)
			break
		}

		msgHead, err := dp.Unpack(byteHead)
		if err != nil {
			fmt.Println(" msg head unpack error ", err)
			break
		}
		msg := msgHead.(*znet.Message)

		if msg.GetDataLen() > 0 {

			msg.Data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println(" read msg body error ", err)
				break
			}
			fmt.Println(" recv msg id: ", msg.Id, " msg len :", msg.GetDataLen(), " msg body: ", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}

}
