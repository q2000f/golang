package network

import (
	"../util"
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"errors"
	"encoding/json"
)

type CallInfo struct {
	Data []byte
	Uuid string
	WaitChan   chan []byte
}

type RpcInfoMap map[string] chan []byte
type RpcCallMgr struct {
	MapRpcInfo RpcInfoMap
	locker sync.Mutex
}

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

func (c *RpcCallMgr) InsertCall(callID string, wait chan []byte) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.MapRpcInfo[callID] = wait
}

func (c *RpcCallMgr) OnCall(callID string, data []byte) {
	c.locker.Lock()
	defer c.locker.Unlock()
	wait, ok := c.MapRpcInfo[callID]
	if ok {
		delete(c.MapRpcInfo, callID)
	}
	wait <- data
}

type RpcConn struct {
	Type           int
	ID             int
	LastPing       int64
	Conn           *net.TCPConn
	reader         *bufio.Reader
	writeChan      chan CallInfo
	closeWaitGroup sync.WaitGroup
	Rpc			   RpcCallMgr
}



func NewServerConn(conn *net.TCPConn) *RpcConn {
	tcpConn := &RpcConn{
			Conn:     conn,
			LastPing: time.Now().Unix(),
	}
	tcpConn.Init()
	go tcpConn.ReadLoop()
	go tcpConn.WriteLoop()
	return tcpConn
}

func ConnectServer(id int, addr string) (*RpcConn, error) {
	tcpAddr, err1 := net.ResolveTCPAddr("tcp", addr)
	if err1 != nil {
		log.Println("ConnectServer net.ResolveTCPAddr", err1, addr)
		return nil, err1
	}

	conn, err2 := net.DialTCP("tcp", nil, tcpAddr)
	if err2 != nil {
		log.Println("ConnectServer DialTCP:", err2, addr)
		return nil, err2
	}
	return NewServerConn(conn), nil
}

func (c *RpcConn) Init() {
	c.reader = bufio.NewReader(c.Conn)
	c.writeChan = make(chan CallInfo, WRITE_CHAN_SIZE)
	c.Rpc.MapRpcInfo = make(RpcInfoMap)
}

func (c *RpcConn) GetID() int {
	return c.ID
}

func (c *RpcConn) Close() {
	c.Conn.Close()
	close(c.writeChan)
}

func (c *RpcConn) WaitClose() {
	c.closeWaitGroup.Wait()
}

func (c *RpcConn) ReadLoop() {
	// read length
	var size uint32
	var sizeBytes = make([]byte, HEADER_SIZE)
	c.closeWaitGroup.Add(1)
	for {
		if _, err := io.ReadFull(c.reader, sizeBytes); err != nil {
			log.Println(err)
			break
		}
		if size = binary.BigEndian.Uint32(sizeBytes); size > MAX_MESSAGE_SIZE {
			log.Println("error package size:", size)
			break
		}
		buff := make([]byte, size)

		if _, err := io.ReadFull(c.reader, buff[:]); err != nil {
			log.Println(err)
			break
		}
		c.OnMessage(buff)
	}
	c.OnClose()
	c.closeWaitGroup.Done()
}

func (c *RpcConn) WriteLoop() {
	c.closeWaitGroup.Add(1)
	defer c.closeWaitGroup.Done()
	sec:= time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sec.C:
		case callInfo, ok := <-c.writeChan:
			if !ok {
				log.Println("WriteLoop break")
				return
			}
			c.Rpc.InsertCall(callInfo.Uuid, callInfo.WaitChan)
			bytes := callInfo.Data
			var sendBytes = make([]byte, HEADER_SIZE + len(bytes))
			binary.BigEndian.PutUint32(sendBytes[0:HEADER_SIZE], uint32(len(bytes)))
			copy(sendBytes[HEADER_SIZE:], bytes)
			err := c.OnWrite(sendBytes)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *RpcConn) Call(opCode int, req interface{}, response interface{}) error {
	out := PacketMessage{
		MessageHeader:MessageHeader{
			OpCode:opCode,
			CallID:util.NewUUID(),
		},
		Body:req,
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		log.Println(bytes)
	}
	wait := make(chan []byte)
	in := CallInfo{
		Data:bytes,
		Uuid:out.CallID,
		WaitChan:wait,
	}
	c.writeChan <- in
	result, ok := <- in.WaitChan
	if !ok {
		return errors.New("rpc chan close!")
	}
	err = json.Unmarshal(result, response)
	if err != nil {
		return err
	}
	return nil
}

func (c *RpcConn) OnMessage(bytes []byte) {
	log.Println("OnMessage", string(bytes))
	message := Message{}
	err := json.Unmarshal(bytes, & message)
	if err != nil {
		log.Println(err)
		return
	}
	c.Rpc.OnCall(message.CallID, message.Body)
}

func (c *RpcConn) OnClose() {
	log.Println("OnClose!")
}

func (c *RpcConn) OnWrite(bytes []byte) error {
	sendSize, err := c.Conn.Write(bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	if sendSize != len(bytes) {
		log.Println("error send len!")
	}
	return nil
}
