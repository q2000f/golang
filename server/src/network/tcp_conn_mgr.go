package network

import (
	"sync"
)

type TcpConnMap map[int] *TcpConn
type TcpConnMgr struct {
	locker  sync.Mutex
	ConnMap TcpConnMap
}

func GetNewTcpConnMgr() *TcpConnMgr {
	tcpConnMgr := &TcpConnMgr{}
	tcpConnMgr.Init()
	return tcpConnMgr
}

func (c *TcpConnMgr) Init() {
	c.ConnMap = TcpConnMap{}
}

func (c *TcpConnMgr) Set(conn *TcpConn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.ConnMap[conn.ID] = conn
}

func (c *TcpConnMgr) Get(id int) (*TcpConn, bool) {
	c.locker.Lock()
	defer c.locker.Unlock()

	conn, ok := c.ConnMap[id]
	return conn, ok
}

func (c *TcpConnMgr) Delete(id int) {
	c.locker.Lock()
	defer c.locker.Unlock()
    _, ok := c.ConnMap[id]
    if ok {
        delete(c.ConnMap, id)
    }
}

func (c *TcpConnMgr) Visit(handFunc func(*TcpConn)) {
    c.locker.Lock()
    defer c.locker.Unlock()
    for _, conn := range c.ConnMap {
        handFunc(conn)
    }
}

func (c *TcpConnMgr) CloseAll() {
    c.locker.Lock()
    defer c.locker.Unlock()
    for _, conn := range c.ConnMap {
        conn.Close()
    }
    c.ConnMap = TcpConnMap{}
}
