package ziface

/*
 消息管理抽象层
*/

type IMsgHandle interface {
	//调度对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	//启动Worker工作池
	StartWorkerPool()
	//消息发送处理
	SendMsgToTaskQueue(request IRequest)
}
