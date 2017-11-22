package client

import (
    "log"
    "../../network"
)

func RegisterServer(conn *network.TcpConn, body []byte) {
    log.Println("RegisterServer:", string(body))
    conn.Status = network.ConnRegister
}

func NotifyRegisterServer(conn *network.TcpConn, body []byte) {
    log.Println("NotifyRegisterServer:", string(body))
}