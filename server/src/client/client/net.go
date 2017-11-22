package client

import (
	"../../network"
	"../../protocol"
	"encoding/json"
	"log"
	"time"
	"../../util"
)

type Client struct {
	conn *network.TcpConn
}

func NewClient(host string) *Client {
	handler := NewServerHandler()
	client := &Client{}
	var err error
	client.conn, err = network.Connect(1, host, handler)
	if err != nil {
		log.Println(err)
	}
	go client.LoopPing()
	return client
}

func (c *Client) SendPacket(code int, in interface{}) {
	out := protocol.PacketMessage{
		MessageHeader: protocol.MessageHeader{
			OpCode: code,
			CallID: util.NewUUID(),
		},
		Body: in,
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		log.Println(bytes)
	}
	log.Println(string(bytes))
	c.conn.SendPacket(bytes)
}

func (c *Client) Ping(in protocol.Ping) {
	c.SendPacket(protocol.PING, in)
}

func (c *Client) RegisterServer(in protocol.IRegisterServer) {
	c.SendPacket(protocol.REG_SERVER, in)
}

func (c *Client) LoopPing() {
	for c.conn.Status != network.ConnClose {
		in := protocol.Ping{
			ClientTime: time.Now().Unix(),
		}
		c.Ping(in)
		time.Sleep(time.Second * 5)
	}
}
func InitNetWork() {
	// rpc, _ := network.ConnectServer(1, "localhost:1236")
	// in := protocol.Ping{
	// 	ClientTime:time.Now().Unix(),
	// }
	// out := protocol.Pong{}
	// rpc.Call(protocol.PING, in, &out)
	// log.Println(out)

	for i := 0; i < 1; i++ {
		client := NewClient("localhost:1236")
		reg := protocol.IRegisterServer{}
		reg.ServerID = 1001 + i;
		reg.ServerType = 1
		reg.ServerName = "game"
		client.RegisterServer(reg)
	}
}
