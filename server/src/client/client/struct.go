package client
import(
	"encoding/json"
	"../../network"
)

type MessageHeader struct {
	OpCode int `json:"opcode"`
}

//Message protocal
type Message struct {
	MessageHeader
	Body   json.RawMessage
}

type HandlerFunc func(conn *network.TcpConn, body []byte)

