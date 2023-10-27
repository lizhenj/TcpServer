package core

import (
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"sync"
	"zinx/app_demo/mmo_game/pb"
	"zinx/ziface"
)

// 玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的链接（用于和客户端的链接）
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标(注意不是Y)
	V    float32            //旋转的0-360角度
}

/*
Player ID 生成器
*/
var PidGen int32 = 1  //玩家ID计数器
var IdLock sync.Mutex //保护PidGen的Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(50)), //随机在160坐标点，基于x轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(50)), //随机在160坐标点，基于Y轴若干偏移
		V:    0,
	}

	return p
}

/*
提供一向客户端发送消息的方法
将pb的protobuf数据序列化后，再调用zinx的sendMsg
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化，转化为二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println("marshal msg err:", err)
		return
	}

	//将msg通过zinx框架sendmsg发给客户端
	if p.Conn == nil {
		log.Println("connection is player is nil")
		return
	}

	if err = p.Conn.SendMsg(msgId, msg); err != nil {
		log.Println("Player SendMsg error!")
	}
}

// 告知客户端玩家Pid
func (p *Player) SyncPid() {
	//组建MsgID:0 的proto数据
	proto_data := &pb.SyncPid{
		Pid: p.Pid,
	}

	//将消息发送客户端
	p.SendMsg(1, proto_data)
}

// 广播玩家字节的初始位置
func (p *Player) BroadCastStartPosition() {
	//组件MsgID:200 的proto数据
	proto_data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //tp:2代表广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将消息发送客户端
	p.SendMsg(200, proto_data)
}

// 玩家广播世界聊天
func (p *Player) Talk(context string) {
	//组件msgID:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: context,
		},
	}
	//得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//向所有玩家发送聊天数据
	for _, player := range players {
		//分别给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}
}

//玩家上线，同步广播玩家位置消息
func (p *Player) SyncSurrounding() {
	//1.获取当前周围玩家（九宫格）
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player, ok := WorldMgrObj.Players[int32(pid)]
		if ok {
			players = append(players, player)
		}
	}

	//2将当前玩家位置信息通过msgID:200 发往周围玩家
	//2.1组件msgID:200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//2.2向周围玩家发送消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//3 获取周围玩家位置信息
	//周围玩家信息切片
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p_msg := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p_msg)
	}

	//封装SyncPlayer protobuf数据
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	p.SendMsg(202, SyncPlayers_proto_msg)
}

func (p *Player) UpdataPos(x, y, z, v float32) {

	//触发消失视野和添加视野业务
	//计算旧格子gid
	oldGid := WorldMgrObj.AoiMgr.GetGidByPos(p.X, p.Z)
	//计算新格子gid
	newGid := WorldMgrObj.AoiMgr.GetGidByPos(x, z)

	//更新玩家位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	if oldGid != newGid {
		//触发gird切换
		//把pid从旧的aoi格子中删除
		WorldMgrObj.AoiMgr.RemovePidFromGrid(int(p.Pid), oldGid)
		//把pid添加到新的aoi格子中去
		WorldMgrObj.AoiMgr.AddPidToGrid(int(p.Pid), newGid)

		p.OnExchangeAoiGrid(oldGid, newGid)
	}
	//组建广播proto协议
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, //移动后的坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}

	//获取当前玩家周边玩家AOI九宫格内的玩家
	players := p.GetSuroundingPlayers()
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) OnExchangeAoiGrid(oldGid, newGid int) error {
	//获取旧的九宫格成员
	oldGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(oldGid)

	//为旧的九宫格成员建立哈希表，用来快速查找
	oldGridsMap := make(map[int]bool, len(oldGrids))
	for _, grid := range oldGrids {
		oldGridsMap[grid.GID] = true
	}

	//获取新的九宫格成员
	newGrids := WorldMgrObj.AoiMgr.GetSurroundGridsByGid(newGid)
	//为新的九宫格成员建立哈希表，快速查找
	newGridsMap := make(map[int]bool, len(newGrids))
	for _, grid := range newGrids {
		newGridsMap[grid.GID] = true
	}

	//-----> 处理视野消失
	offline_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//找到在旧的九宫格中出现，但是在新的九宫格没出现的格子
	leavingGrids := make([]*Grid, 0)
	for _, grid := range oldGrids {
		if _, ok := newGridsMap[grid.GID]; !ok {
			leavingGrids = append(leavingGrids, grid)
		}
	}

	//获取需要消失的格子中的全部玩家
	for _, grid := range leavingGrids {
		players := WorldMgrObj.GetPlayersByGid(grid.GID)
		for _, player := range players {
			//让自己在其他玩家的客户端消失
			player.SendMsg(201, offline_msg)

			//将其他玩家信息，在自己客户端消失
			another_offline_msg := &pb.SyncPid{
				Pid: player.Pid,
			}
			p.SendMsg(201, another_offline_msg)
		}
	}

	//------》处理视野出现
	//找到在新九宫格中出现，旧的没有的格子
	enteringGrids := make([]*Grid, 0)
	for _, grid := range newGrids {
		if _, ok := oldGridsMap[grid.GID]; !ok {
			enteringGrids = append(enteringGrids, grid)
		}
	}

	online_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取需要显示格子的全部玩家
	for _, grid := range enteringGrids {
		players := WorldMgrObj.GetPlayersByGid(grid.GID)

		for _, player := range players {
			//让自己出现在其他人视野中
			player.SendMsg(200, online_msg)

			//让其他人出现在自己视野中
			another_online_msg := &pb.BroadCast{
				Pid: player.Pid,
				Tp:  2,
				Data: &pb.BroadCast_P{
					&pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			p.SendMsg(200, another_online_msg)
		}
	}

	return nil
}

func (p *Player) GetSuroundingPlayers() []*Player {
	//得到当前AOI九宫格内所有玩家ID
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player, ok := WorldMgrObj.Players[int32(pid)]
		if ok {
			players = append(players, player)
		}
	}

	return players
}

func (p *Player) Offline() {
	//得到当前玩家周边九宫格玩家
	players := p.GetSuroundingPlayers()

	//封装消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
