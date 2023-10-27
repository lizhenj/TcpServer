package main

import (
	"log"
	"zinx/app_demo/mmo_game/apis"
	"zinx/app_demo/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

// 客户端建立链接和的hook函数
func onConnectionAdd(conn ziface.IConnection) {
	//创建player对象
	player := core.NewPlayer(conn)

	//给客户端发送MsgID:1的消息
	player.SyncPid()
	//给客户端发送MsgID:200的消息
	player.BroadCastStartPosition()

	//将新上线玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个pid 玩家id
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家，告知新上线玩家位置信息
	player.SyncSurrounding()

	log.Println("====> Player pid=", &player.Pid, " is arrived<===")
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, _ := conn.GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	player.Offline()

	log.Println("=====> Player pid = ", pid, " offline...<==")
}

func main() {
	//创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//链接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(onConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Serve()
}
