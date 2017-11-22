package consul

import (
	"../../protocol"
	"../../network"
	"encoding/json"
	"log"
	"net"
	"time"
	"sync"
)

type HandlerFunc func(conn *network.TcpConn, body []byte) (interface{}, error)
type TcpConnMap map[int]*network.TcpConn
type HandlerFuncMap map[int]HandlerFunc

type ServerHandler struct {
	MaxID         int
	handlerMap    map[int]HandlerFunc
	ClientConnMap TcpConnMap
	mapLock sync.RWMutex
}
//NewServerHandler create handler
func NewServerHandler(handlerMap HandlerFuncMap) *ServerHandler {
	server := &ServerHandler{handlerMap:handlerMap}
	server.Init()
	go server.Run()
	return server
}

//Init map
func (c *ServerHandler) Init() {
	c.ClientConnMap = make(TcpConnMap)
}
//CheckPing is check ping
func (c *ServerHandler) CheckPing() {
	c.mapLock.RLock()
	defer c.mapLock.RUnlock()
	conns := []*network.TcpConn{}
	for _, conn := range c.ClientConnMap {
		now := time.Now().Unix()
		if now - conn.LastPing > PING_TIME_LIMIT {
			conns = append(conns, conn)
		}
	}

	for _, conn := range conns {
		conn.Close()
	}
}
//BroadCast is broadcast to all connted socket.
func (c *ServerHandler) BroadCast(opcode int, data interface{}) {
	c.mapLock.RLock()
	defer c.mapLock.RUnlock()
	for _, conn := range c.ClientConnMap {
		c.SendPacket(conn, opcode, data)
	}
}

func (c *ServerHandler) Run() {
	second := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-second.C:
			c.CheckPing()
		}
	}
}

func (c *ServerHandler) OnConnect(conn *net.TCPConn) {
	newConn := network.NewTcpConn(conn, c)	
	log.Println("OnConnect ", "ID:", newConn.GetID(), " addr:", conn.RemoteAddr())
	c.mapLock.Lock()
	defer c.mapLock.Unlock()
	c.ClientConnMap[newConn.GetID()] = newConn
}

func (c *ServerHandler) OnMessage(conn *network.TcpConn, bytes []byte) {
	log.Println("OnMessage", string(bytes))
	conn.LastPing = time.Now().Unix()
	message := protocol.Message{}
	err := json.Unmarshal(bytes, &message)
	if err != nil {
		log.Println(err)
		return
	}

	handlerFunc, ok := c.handlerMap[message.OpCode]
	if !ok {
		log.Println("error opcod:", message)
		return
	}
	out, err := handlerFunc(conn, message.Body)
	if err != nil {
		log.Println(err)
		return
	}
	c.SendCallBack(conn, message.CallID, message.OpCode, out)
}

func (c *ServerHandler) OnClose(conn *network.TcpConn) {
	log.Println("OnClose", "ID:", conn.GetID())
	c.mapLock.Lock()
	defer c.mapLock.Unlock()
	delete(c.ClientConnMap, conn.GetID())	
	ServerMgrInstance.OffLine(conn.AttachID)
}

func (c *ServerHandler) GetNewID() int {
	c.MaxID++
	return c.MaxID
}

func (c *ServerHandler) SendCallBack(conn *network.TcpConn, callID string, code int, data interface{}) error {
	if data  == nil {
		return nil
	}
	packet := protocol.PacketMessage{
		MessageHeader:protocol.MessageHeader{
			OpCode:code,
			CallID:callID,
		},
		Body:data,
	}

	bytes, err := protocol.Marshal(packet)
	if err != nil {
		return err
	}
		log.Println("SendCallBack",string(bytes))
	conn.SendPacket(bytes)
	return nil
}

func (c *ServerHandler) SendPacket(conn *network.TcpConn, code int, data interface{}) error {
	c.SendCallBack(conn, "", code, data)
	return nil
}