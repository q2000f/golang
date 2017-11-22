package consul

import (
	"../../network"
)

//InitNetWork net
var ServerConnMgr *ServerHandler
func InitNetWork() {
	ServerConnMgr = NewServerHandler(InitServerHandler())
	network.NewAccept("localhost:1236", ServerConnMgr)
}
