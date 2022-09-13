package znet

import (
	"fmt"
	"test/iface"
	"test/utils"
)

type Server struct {
	ServerName string
	IpVersion  string
	Ip         string
	Port       int

	MsgHandler iface.IMessageHandler

	//改用kqueue以后这个没用了
	//	ConnManager iface.IConnectionManager

	OnEventLoopStart func(el iface.IEventLoop)

	OnEventLoopStop func(el iface.IEventLoop)

	ELoop iface.IEventLoop
}

//这块是不用的，现在用kqueue作为reactor不需要这个启动了
/*func (s *Server) Start() {
	go func() {
		fmt.Println("[start] server base version is starting...")
		s.MsgHandler.StartWorkerPool()
		fmt.Printf("[Start] server listener at ip :%s, port %d, is starting /n", s.Ip, s.Port)
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IpVersion, " error: ", err)
			return
		}

		fmt.Println("start server succ, ", s.ServerName, " listening... ")

		var cid uint32
		cid = 0

		for {
			var conn, err = listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept error: ", err)
				continue
			}

			if s.ConnManager.GetConnCount() >= utils.GlobalObject.MaxConn {
				fmt.Println(" too many connection ")
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid += 1

			go dealConn.Start()

		}
	}()

}*/

func (s *Server) StartKqueueServer() {
	go func() {
		fmt.Println("[start] Kserver kqueue version is starting... ")
		s.MsgHandler.StartWorkerPool()
		fmt.Printf("[Start] Kserver listener at ip :%s, port %d, is starting ", s.Ip, s.Port)
		fmt.Println("  ")
		serverSockerUnit, err := Listen(s.Ip, s.Port)
		if err != nil {
			fmt.Println("Kserver init error, ", err)
			return
		}
		eventLoop, err := NewEventLoop(serverSockerUnit, s)
		s.ELoop = eventLoop
		if err != nil {
			fmt.Println("Kserver new eventLoop error, ", err)
			return
		}
		fmt.Println("Kserver start looping ")
		eventLoop.Looping()

	}()
	s.CallOnEventLoopStart(s.ELoop)
}

func (s *Server) Stop() {
	fmt.Println(" server name: ", s.ServerName, " stopped! ")

	s.CallOnEventLoopStop(s.ELoop)
	s.ELoop.Stop()

}

func (s *Server) Serve() {
	//	s.Start()
	s.StartKqueueServer()
	select {}
}

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println(" add router to ", s.ServerName, " succ ")

}

func NewServer() iface.IServer {
	s := &Server{
		ServerName: utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMessageHandler(),
		//ConnManager: NewConnectionManager(),
	}
	return s
}

// 这块原来是conn用的，现在用socketunit代替了
/*func (s *Server) GetConnectionManager() iface.IConnectionManager {
	return s.ConnManager

}*/

func (s *Server) SetOnSocketUnitStart(startHookFunc func(su iface.ISocketUnit)) {
	s.ELoop.SetOnSocketUnitStart(startHookFunc)
}

func (s *Server) SetOnSocketUnitStop(stopHookFunc func(su iface.ISocketUnit)) {
	s.ELoop.SetOnSocketUnitStop(stopHookFunc)
}

func (s *Server) CallOnEventLoopStart(el iface.IEventLoop) {
	if s.OnEventLoopStop != nil {
		fmt.Println(" call onConnStart ")
		s.OnEventLoopStop(el)
	}
}

func (s *Server) CallOnEventLoopStop(el iface.IEventLoop) {
	if s.OnEventLoopStop != nil {
		fmt.Println(" call onConnStop")
		s.OnEventLoopStop(el)
	}
}

func (s *Server) SetOnEventLoopStart(startHookFunc func(el iface.IEventLoop)) {
	s.OnEventLoopStart = startHookFunc
}
func (s *Server) SetOnEventLoopStop(stopHookFunc func(el iface.IEventLoop)) {
	s.OnEventLoopStop = stopHookFunc
}
