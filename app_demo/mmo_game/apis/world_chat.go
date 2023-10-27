package apis

import (
	"google.golang.org/protobuf/proto"
	"log"
	"zinx/app_demo/mmo_game/core"
	"zinx/app_demo/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

// 世界聊天路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(req ziface.IRequest) {
	//1解析客户端传递来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), proto_msg)
	if err != nil {
		log.Println("Talk Unmarshal error ", err)
		return
	}

	//当前聊天数据属于哪个玩家发送
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		log.Println("Talk Unmarshal error ", err)
		return
	}

	//3 根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4 将聊天数据广播其他在线的玩家
	player.Talk(proto_msg.Content)
}
