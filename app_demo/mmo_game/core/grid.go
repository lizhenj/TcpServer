package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int          //格子标号
	MinX      int          //格子左边边界坐标
	MaxX      int          //格子右边边界坐标
	MinY      int          //格子上边边界坐标
	MaxY      int          //格子下边边界坐标
	playerIDs map[int]bool //当前格子内玩家或物体成员的ID集合
	pIDLock   sync.RWMutex
}

// 初始化当前格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 格子添加玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 格子删除玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 当前格子的玩家ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, "+
		"playerIDs:%v\n", g.GID, g.MinX, g.MaxY, g.MinY, g.MaxX, g.playerIDs)
}
