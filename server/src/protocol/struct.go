package protocol

import "encoding/json"

//MessageHeader header
type MessageHeader struct {
	OpCode int    `json:"code"`
	CallID string `json:"cid"` //call id
}

//Message protocal
type Message struct {
	MessageHeader
	Body json.RawMessage
}

type PacketMessage struct {
	MessageHeader
	Body interface{}
}

func PackNotify(opcode int, in interface{}) ([]byte, error) {
	message := PacketMessage{}
	message.OpCode = opcode
	message.Body = in
	bytes, err := json.Marshal(message)
	return bytes, err
}

func Pack(opcode int, callID string, in interface{}) ([]byte, error) {
	message := PacketMessage{}
	message.OpCode = opcode
	message.CallID = callID
	message.Body = in
	bytes, err := json.Marshal(message)
	return bytes, err
}

func Marshal(in interface{}) ([]byte, error) {
	bytes, err := json.Marshal(in)
	return bytes, err
}

func Unmarshal(out interface{}, data []byte) error {
	err := json.Unmarshal(data, out)
	return err
}

///////////////////////////////
type Ping struct {
	ClientTime int64 `json:"ct"`
}

type Pong struct {
	ClientTime int64 `json:"ct"`
	ServerTime int64 `json:"st"`
}

const (
	ServerRegister = 1
	ServerUpdate   = 2
	ServerRemove   = 3
)
type ServerInfo struct {
	ServerID   int
	ServerType int
	ServerName string
}

type IRegisterServer struct {
	ServerInfo
}

type ORegisterServer struct {
}

type NotifyRegisterServer struct {
	ServerInfo
	Action int
}

type IGetServerList struct {

}

type OGetServerList struct {
	ServerList []ServerInfo
}
