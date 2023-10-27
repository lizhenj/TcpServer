package apis

import (
	"google.golang.org/protobuf/proto"
	"log"
	"zinx/app_demo/mmo_game/core"
	"zinx/app_demo/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

/*
 玩家移动
*/
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(req ziface.IRequest) {
	//解析客户端传来proto协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), proto_msg)
	if err != nil {
		log.Println("Move: Position Unmarshal error ", err)
		return
	}

	//得到当前发送位置的玩家
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		log.Println("GetProperty pid error ", err)
		return
	}
	log.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", pid,
		proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//给其他玩家进行当前玩家的位置广播信息
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	if player == nil {
		return
	}

	player.UpdataPos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

}
