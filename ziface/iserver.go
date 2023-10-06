package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//路由功能，当前的服务注册一路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)

	//获取当前server的链接管理器
	GetConnMgr() IConnManager

	//注册及调用onConnStart、stop钩子
	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
