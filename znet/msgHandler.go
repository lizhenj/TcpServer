package znet

import (
	"log"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
	"zinx/zlog"
)

/*
 消息处理模块实现
*/

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

//初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 获取msgID对应的handler
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		log.Printf("api msgID=%d is NOT FOUND!\n", request.GetMsgID())
		return
	}
	//调度对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 对当前msg是否已存在绑定API处理方法进行判断
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api , msgID=" + strconv.Itoa(int(msgID)))
	}
	//绑定msg与API
	mh.Apis[msgID] = router
	log.Printf("Add api MsgID=%d SUCC!\n", msgID)
}

//启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize,分别开启Worker,每个Worker用新协程承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//启动worker
		//1.为当前worker对应的channel消息队列开辟空间（缓存通道）
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前worker，阻塞等待channel消息
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	zlog.Infof("Worker ID=%d is started...", workerID)

	//不断的阻塞等待对应消息队列的消息,并进行处理
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息发往TaskQueue,由Worker处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同通道的worker
	workerID := request.GetMsgID() % mh.WorkerPoolSize
	log.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(),
		" to WorkerID = ", workerID)

	mh.TaskQueue[workerID] <- request
}
