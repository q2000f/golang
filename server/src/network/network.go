package network

import (
	"net"
	"log"
)

type TcpHandler interface {
	OnConnect(conn *TcpConn)
	OnMessage(conn *TcpConn, bytes []byte)
	OnClose(conn *TcpConn)
}

type NetServer struct {
	listen    net.Listener
	tcpConnMgr   *TcpConnMgr
	tcpHandler TcpHandler
	ip string
	id int
}

func StartServer(ip string, handler TcpHandler) *NetServer {
	server := &NetServer{
		ip:ip,
		tcpHandler:handler,
		tcpConnMgr:GetNewTcpConnMgr(),
	}
	server.Init()
	return server
}

func (c *NetServer) Init() error {
	var err error
	c.listen, err = net.Listen("tcp", c.ip)
	if err != nil {
		log.Println("StartListen net.ListenTCP err:", err)
		return err
	}
	log.Println("listen on:", c.ip)
	go c.run()
	return nil
}

func (c *NetServer) run() {
	for {
		conn, err := c.listen.Accept()
		if err != nil {
			log.Println("listen.Accept err:", conn)
			break
		}
		c.onConnect(conn.(*net.TCPConn))
	}
}
func (c *NetServer) onConnect(conn *net.TCPConn) {
	tcpConn := NewTcpConn(conn, c.GetNewID(), c)
	c.tcpConnMgr.Set(tcpConn)
	c.tcpHandler.OnConnect(tcpConn)
}

func (c *NetServer) OnMessage(conn *TcpConn, bytes []byte) {
	c.tcpHandler.OnMessage(conn, bytes)
}

func (c *NetServer) OnClose(conn *TcpConn) {
	c.tcpHandler.OnClose(conn)
	c.tcpConnMgr.Delete(conn.GetID())
}

func (c *NetServer) Close(id int) {
	conn, ok := c.tcpConnMgr.Get(id)
	if ok {
		conn.Close()
	}
}

func (c *NetServer) WaitClose(id int) {
	conn, ok := c.tcpConnMgr.Get(id)
	if ok {
		conn.WaitClose()
	}
}

func (c *NetServer) CloseAll() {
	c.tcpConnMgr.Visit( func(tcpConn *TcpConn) {
		tcpConn.Close()
	})
}

func (c *NetServer) WaitAllClose() {
	c.tcpConnMgr.Visit( func(tcpConn *TcpConn) {
		tcpConn.WaitClose()
	})
}

func (c *NetServer) GetNewID() int {
	c.id++;
	return c.id
}