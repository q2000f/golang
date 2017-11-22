package server

import "../../network"

var server *network.NetServer
func Start() {
    InitServerMgr()
    server = network.StartServer("localhost:1236", NewTcpHandler())
}

func Stop() {
    server.CloseAll()
}
