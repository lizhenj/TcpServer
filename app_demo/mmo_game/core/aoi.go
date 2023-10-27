package core

import (
	"fmt"
)

// 定义一些AOI的边界值
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
 AOI区域管理模块
*/

type AOIManager struct {
	MinX  int           //区域的左边界坐标
	MaxX  int           //区域的右边界坐标
	CntsX int           //X方向的格子数量
	MinY  int           //区域的上边界坐标
	MaxY  int           //区域的下边界坐标
	CntsY int           //Y方向的格子数量
	grids map[int]*Grid //当前区域中的格子
}

// AOI区域管理初始化
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		minX,
		maxX,
		cntsX,
		minY,
		maxY,
		cntsY,
		make(map[int]*Grid),
	}

	//AOI区域初始化所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//根据y与x计算坐标编号
			//id = idy*cntX + idx
			gid := y*cntsX + x

			//初始化gid格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

// 获取每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 获取每个格子在y轴方向的高度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, cntsX:%d, "+
		"minY:%d, maxY:%d, cntsY:%d\n Grids in AOIManager:\n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	for _, grid := range m.grids {
		s += fmt.Sprintf("%s", grid)
	}

	return s
}

// 根据格子GID得到周边九宫格格子的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gID是否在aoi中
	grid, ok := m.grids[gID]
	if !ok {
		return
	}

	//初始化grids返回值切片
	grids = append(grids, grid)

	//gID的左右边是否存在格子
	// idx = id % nx
	idx := gID % m.CntsX

	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	//遍历已获得的X轴，查看其上下格子的存在与否
	for _, v := range gidsX {
		//idy = id / nx
		idy := v / m.CntsX
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

// 通过x,y坐标获得GID格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部的PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//获得当前玩家的GID格子id
	gID := m.GetGidByPos(x, y)

	//获取九宫格范围值
	grids := m.GetSurroundGridsByGid(gID)

	//获取范围内所有玩家id
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		//log.Printf("==> grid ID:%d, pids: %v ====", v.GID, v.GetPlayerIDs())
	}
	return
}

// 添加playerID到格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	if grid, ok := m.grids[gID]; ok {
		grid.Add(pID)
	}
}

// 移除格子中的playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	if grid, ok := m.grids[gID]; ok {
		grid.Remove(pID)
	}
}

// 通过GID获取全部的playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	grid, ok := m.grids[gID]
	if !ok {
		return
	}
	playerIDs = grid.GetPlayerIDs()
	return
}

// 通过坐标将player添加到格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid, ok := m.grids[gID]
	if !ok {
		return
	}
	grid.Add(pID)
}

// 通过坐标将player从格子删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid, ok := m.grids[gID]
	if !ok {
		return
	}
	grid.Remove(pID)
}
