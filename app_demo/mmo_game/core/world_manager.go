package core

import "sync"

/*
 当前游戏的世界总管理模块
*/

type WorldManager struct {
	//AOIManager 当前世界地图AOI的管理模块
	AoiMgr *AOIManager
	//当前在线的players集合
	Players map[int32]*Player
	//保护Players集合的锁
	pLock sync.RWMutex
}

// 提供一个对外的世界管理模块句柄（全局）
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		//创建世界AOI地图规划
		AoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X,
			AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		//初始化player集合
		Players: make(map[int32]*Player),
	}
}

// 添加玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player添加到AOIManager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// 删除玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	//得到玩家
	wm.pLock.Lock()
	defer wm.pLock.Unlock()

	player, ok := wm.Players[pid]
	if !ok {
		return
	}
	wm.AoiMgr.RemoveFromGridByPos(int(player.Pid), player.X, player.Z)

	delete(wm.Players, pid)
}

// 查询玩家
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

//获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, player := range wm.Players {
		players = append(players, player)
	}

	return players
}

//获取指定gid中的所有player信息
func (wm *WorldManager) GetPlayersByGid(gid int) []*Player {
	//通过gid获取对应格子中的所有pid
	pids := wm.AoiMgr.grids[gid].GetPlayerIDs()

	//通过pid找到对应的player对象
	players := make([]*Player, 0, len(pids))
	wm.pLock.RLock()
	for _, pid := range pids {
		players = append(players, wm.Players[int32(pid)])
	}
	wm.pLock.RUnlock()

	return players
}