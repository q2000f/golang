package consul

import (
)
import (
    "sync"
    "../../protocol"
    "../../network"
)

type ServerInfo struct {
    protocol.ServerInfo
    Conn *network.TcpConn
}

type MapServerInfo map[int]ServerInfo

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

    c.ServerInfoMap[serverInfo.ServerID] = serverInfo

    notify := protocol.NotifyRegisterServer{}
    notify.ServerInfo = serverInfo.ServerInfo
    notify.Action = protocol.ServerRegister
    
    ServerConnMgr.BroadCast(protocol.NOTIFY_REG_SERVER, notify)
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

func (c *ServerMgr) OffLine(serverID int) {
    c.locker.Lock()
    defer c.locker.Unlock()
    delete(c.ServerInfoMap, serverID)
}

