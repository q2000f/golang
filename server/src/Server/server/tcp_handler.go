package server

import (
    "../../network"
    "log"
    "protocol"
    "encoding/json"
)

type HandlerFunc func(conn *network.TcpConn, body []byte) (interface{}, error)
type TcpHandler struct{
    handlerMap map[int]HandlerFunc
}
//
func NewTcpHandler() *TcpHandler {
    tcpHandler := &TcpHandler{}
    tcpHandler.Init()
    return tcpHandler
}

func (c *TcpHandler) Init() {
    c.handlerMap = make(map[int]HandlerFunc)
    c.handlerMap[protocol.PING] = Ping
    c.handlerMap[protocol.REG_SERVER] = RegisterServer
}

func (c *TcpHandler) OnConnect(conn *network.TcpConn) {
    log.Println("OnConnect", conn.ID)
}

func (c *TcpHandler) OnMessage(conn *network.TcpConn, bytes []byte) {
    //log.Println("OnMessage", conn.ID, " ", string(bytes))
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
    out, err := handlerFunc(conn, message.Body)
    if err != nil {
        log.Println(err)
        return
    }
    bytes, errPack := protocol.Pack(message.OpCode, message.CallID, out)
    if errPack != nil {
        log.Println(errPack)
        return
    }
    conn.SendPacket(bytes)
}

func (c *TcpHandler) OnClose(conn *network.TcpConn) {
    log.Println("OnClose", conn.ID)
}
