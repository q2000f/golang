package client

import (
    "log"
    "../../protocol"
    "../../network"
    "net"
    "encoding/json"
)
//
type ServerHandler struct {
    MaxId int
    handlerMap map[int]HandlerFunc
}

//
func NewServerHandler() *ServerHandler {
    server := &ServerHandler{}
    server.Init()
    return server
}

//
func (c *ServerHandler) Init() {
    c.handlerMap = make(map[int]HandlerFunc)
    c.handlerMap[protocol.PING] = Pong
    c.handlerMap[protocol.REG_SERVER] = RegisterServer
    c.handlerMap[protocol.NOTIFY_REG_SERVER] = NotifyRegisterServer
}

func (c *ServerHandler) OnConnect(conn *net.TCPConn) {
    newConn := network.NewTcpConn(conn, c.GetNewID(), c)
    log.Println("OnConnect ", "ID:", newConn.GetID(), " addr:", conn.RemoteAddr())
}

func (c *ServerHandler) OnMessage(conn *network.TcpConn, bytes []byte) {
    log.Println("OnMessage", string(bytes))
    message := protocol.Message{}
    err := json.Unmarshal(bytes, &message)
    if err != nil {
        log.Println(err)
        return
    }

    handlerFunc, ok := c.handlerMap[message.OpCode]
    if !ok {
        log.Println("error opcod:", string(bytes))
        return
    }
    handlerFunc(conn, message.Body)
}

func (c *ServerHandler) OnClose(conn *network.TcpConn) {
    log.Println("OnClose", "ID:", conn.GetID())
}

func (c *ServerHandler) GetNewID() int {
    c.MaxId++
    return c.MaxId
}
