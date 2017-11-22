package consul

import (
    "../../protocol"
)

func InitServerHandler() HandlerFuncMap {
    serverMap := make(HandlerFuncMap)
    serverMap[protocol.PING] = Ping
    serverMap[protocol.REG_SERVER] = RegisterServer
    return serverMap
}