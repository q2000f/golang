package consul

import (
	"../../network"
	"../../protocol"
	"log"
	"errors"
)

//RegisterServer is register server
func RegisterServer(conn *network.TcpConn, body []byte) (interface{}, error) {
    log.Println("RegisterServer:", string(body))
    
    if conn.Status != network.ConnInit {
        result := "has inited!"
        return nil, errors.New(result)
    }

	req := protocol.IRegisterServer{}
	err := protocol.Decode(&req, body)
	res := &protocol.ORegisterServer{}
	if err != nil {
		result := "error!"
		return nil, errors.New(result)
	}
   
	server := ServerInfo{
        ServerInfo:req.ServerInfo,
		Conn:       conn,
	}

	conn.AttachID = req.ServerID
    conn.Status = network.ConnRegister
	ServerMgrInstance.Register(server)
	return res, nil
}
