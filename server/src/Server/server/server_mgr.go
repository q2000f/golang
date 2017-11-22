package server

import (
)
import (
    "sync"
    "../../protocol"
    "../../network"
    "log"
)

type ServerInfo struct {
    protocol.ServerInfo
    Conn *network.TcpConn
}

type MapServerInfo map[int]*ServerInfo

type ServerMgr struct {
    ServerInfoMap MapServerInfo
    locker sync.Mutex
}
var ServerMgrInstance ServerMgr

func InitServerMgr() {
    ServerMgrInstance.ServerInfoMap = MapServerInfo{}
}
func (c *ServerMgr) Register(serverInfo ServerInfo) {
    c.locker.Lock()
    defer c.locker.Unlock()

    oldServer, ok := c.ServerInfoMap[serverInfo.ServerID]
    if ok {
        oldServer.Conn.AttachID = 0
        oldServer.Conn.Close()    
        delete(c.ServerInfoMap, serverInfo.ServerID)
    }
    c.ServerInfoMap[serverInfo.ServerID] = &serverInfo
    notify := protocol.NotifyRegisterServer{}
    notify.ServerInfo = serverInfo.ServerInfo
    notify.Action = protocol.ServerRegister
    packet, err := protocol.PackNotify(protocol.NOTIFY_REG_SERVER, notify)
    if err != nil {
        log.Println(err);
        return
    }
    c.BroadCast(packet)
}

func (c *ServerMgr) GetServerList() []protocol.ServerInfo {
    c.locker.Lock()
    defer c.locker.Unlock()
    serverList := []protocol.ServerInfo{}
    for _, server := range c.ServerInfoMap {
        serverList = append(serverList, server.ServerInfo)
    }
    return serverList
}

//warning no locker
func (c *ServerMgr) BroadCast(notify []byte) {
    for _, serverInfo := range c.ServerInfoMap {
        serverInfo.Conn.SendPacket(notify);
    }
}

func (c *ServerMgr) LockedAndBroadCast(notify []byte) {
    c.locker.Lock()
    defer c.locker.Unlock()
    for _, serverInfo := range c.ServerInfoMap {
        serverInfo.Conn.SendPacket(notify);
    }
}

func (c *ServerMgr) LockedAndSendPacket(serverID int, out []byte) bool {
    c.locker.Lock()
    defer c.locker.Unlock()
    serverInfo, ok := c.ServerInfoMap[serverID]
    if ok {
        serverInfo.Conn.SendPacket(out)
        return true
    }
    return false
}

func (c *ServerMgr) OffLine(serverID int) {
    c.locker.Lock()
    defer c.locker.Unlock()
    delete(c.ServerInfoMap, serverID)
}

