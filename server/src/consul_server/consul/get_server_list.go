package consul

import (
	"../../network"
	"../../protocol"
	"log"
)

//GetServerList is client get server list
func GetServerList(conn *network.TcpConn, body []byte) (interface{}, error) {
    log.Println("GetServerList:", string(body))
    
	// req := protocol.IGetServerList{}
	// err := protocol.Decode(&req, body)
	// 	if err != nil {
	// 	result := []byte("error!")
	// 	return result, nil
	// }
	res := &protocol.OGetServerList{}
	res.ServerList = ServerMgrInstance.GetServerList()
	return res, nil
}
