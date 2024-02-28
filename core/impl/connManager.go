package impl

import (
	"sync"

	"github.com/ikun666/tcp_server/core/iface"
)

// 连接管理器
type connManager struct {
	connMap  map[uint32]iface.IConnection //保存连接
	connLock sync.RWMutex
}

// 添加连接
func (c *connManager) Add(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connMap[conn.GetID()] = conn
}

// 删除连接
func (c *connManager) Remove(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connMap, conn.GetID())
}

// 获取连接
func (c *connManager) Get(id uint32) iface.IConnection {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	return c.connMap[id]
}

// 得到连接数
func (c *connManager) Len() int32 {
	return int32(len(c.connMap))
}

// 清空连接
func (c *connManager) ClearConn() {
	for _, v := range c.connMap {
		//删除连接
		c.Remove(v)
		v.Stop()
	}
}
func newConnManager() iface.IConnManager {
	return &connManager{
		connMap: make(map[uint32]iface.IConnection),
	}
}
