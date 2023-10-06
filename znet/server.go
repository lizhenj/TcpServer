package znet

import (
	"fmt"
	"log"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	////当前Server的router,server注册的链接对应的处理业务
	//Router ziface.IRouter
	//当前server的消息管理模块,与MsgID对应的处理业务API进行绑定
	MsgHandler ziface.IMsgHandle

	//该server的链接管理器
	ConnMgr ziface.IConnManager
	//该server创建链接之后调用的Hook函数-OnConnStart
	OnConnStart func(connection ziface.IConnection)
	//该server销毁链接之前调用的Hook函数-OnConnStop
	OnConnStop func(connection ziface.IConnection)
}

////定义当前客户端链接所绑定的handle api（目前写死，后边用户自定义）
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	//回显的业务
//	log.Println("[Conn Handle] CallBackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		log.Println("write back buf err: ", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Start() {
	log.Printf("[Zinx] Server Name: %s,listenner "+
		"at IP: %s, Port: %d is starting\n", utils.GlobalObject.Name,
		utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	log.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	//log.Printf("[Start] Server Listenner "+
	//	"at IP %s, Port %d, is starting\n", s.IP, s.Port)

	//避免阻塞，另起协程
	go func() {
		//0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Printf("resolve tcp addr error: %v\n", err)
			panic(err)
		}

		// 2监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Printf("listen %s error: %v\n", s.IPVersion, err)
			panic(err)
		}

		log.Println("start Zinx server success, ", s.Name, " success, Listenning...")

		var cid uint32

		// 3阻塞的等待客户端链接，并处理客户端链接业务（读写）
		for {
			//阻塞，直到客户端链接
			conn, err := listenner.AcceptTCP()
			if err != nil {
				log.Printf("Accept err: %v\n", err)
				continue
			}

			//进行最大链接个数的判断，若超出，则关闭此链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//todo 给客户端响应最大链接的错误
				log.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新链接的业务方法与conn进行绑定，得到我们的链接模块对象
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的链接业务处理
			go dealConn.Start()

			////已经与客户端简历链接，进入业务处理。
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			log.Println("recv buf err: ", err)
			//			break
			//		}
			//
			//		//回复
			//		if _, err = conn.Write(buf[:cnt]); err != nil {
			//			log.Println("write back buf err: ", err)
			//			break
			//		}
			//	}
			//}()
		}
	}()
}

func (s *Server) Stop() {
	//todo 将一些服务器的资源、状态或者一些已经开辟的链接信息，进行停止或者回收
	log.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//todo 启动服务器后的额外业务
	//阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
  初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//注册钩子函数及调用
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		log.Println("---> Call OnConnStart()")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		log.Println("---> Call OnConnStop()")
		s.OnConnStop(connection)
	}
}
