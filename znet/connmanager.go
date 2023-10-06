package znet

import (
	"errors"
	"log"
	"sync"
	"zinx/ziface"
)

/*
 链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

//创建当前链接管理的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	//保护资源map。加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
	log.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除链接
	delete(c.connections, conn.GetConnID())
	log.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; !ok {
		return nil, errors.New("connection not FOUND!")
	} else {
		return conn, nil
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn，并停止conn相关工作
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}

	log.Println("Clear All connections succ! conn num = ", c.Len())
}
