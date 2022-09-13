package znet

import (
	"fmt"
	"strconv"
	"test/iface"
	"test/utils"
)

type MessageHandler struct {
	Apis map[uint32]iface.IRouter

	TaskQueues []chan iface.IKRequest

	WorkPoolSize uint32

	MaxTaskQueueLen uint32
}

func (mh *MessageHandler) DoMsgHandle(krequest iface.IKRequest) {
	handler, ok := mh.Apis[krequest.GetMsgId()]
	if ok == false {
		fmt.Println(" handler msgid: ", krequest.GetMsgId(), " is not exist in the map")
		return
	}

	handler.PreHandle(krequest)
	handler.Handle(krequest)
	handler.PostHandle(krequest)
}

func (mh *MessageHandler) AddRouter(msgId uint32, router iface.IRouter) {

	if _, ok := mh.Apis[msgId]; ok {
		panic(" router id :" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println(" add api msgid: ", msgId, " succ! ")
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Apis:            make(map[uint32]iface.IRouter),
		WorkPoolSize:    utils.GlobalObject.WorkPoolSize,
		TaskQueues:      make([]chan iface.IKRequest, utils.GlobalObject.WorkPoolSize),
		MaxTaskQueueLen: utils.GlobalObject.MaxWorkTaskLen,
	}
}

func (mh *MessageHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		mh.TaskQueues[i] = make(chan iface.IKRequest, mh.MaxTaskQueueLen)
		go mh.StartOneWorker(i, mh.TaskQueues[i])
	}

}

func (mh *MessageHandler) StartOneWorker(workerId int, taskQue chan iface.IKRequest) {
	fmt.Println(" worker id: ", workerId, " started... ")
	for {
		select {
		case krequest := <-taskQue:
			mh.DoMsgHandle(krequest)
		}
	}

}

func (mh *MessageHandler) SendRequestToTaskQue(krequest iface.IKRequest) {
	workerId := uint32(krequest.GetSocketUnit().GetSocketFd()) % mh.WorkPoolSize
	fmt.Println(" add socketfd: ", krequest.GetSocketUnit().GetSocketFd(), " msgid: ", krequest.GetMsgId(), " to workerid: ", workerId)
	mh.TaskQueues[workerId] <- krequest
}

