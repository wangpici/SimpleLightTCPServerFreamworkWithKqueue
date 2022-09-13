package iface

type IRouter interface {
	PreHandle(request IKRequest)
	Handle(request IKRequest)
	PostHandle(request IKRequest)
}
