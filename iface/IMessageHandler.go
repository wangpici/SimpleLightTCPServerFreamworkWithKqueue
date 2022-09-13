package iface

type IMessageHandler interface {
	DoMsgHandle(request IKRequest)

	AddRouter(msgId uint32, router IRouter)

	StartWorkerPool()

	SendRequestToTaskQue(request IKRequest)
}
