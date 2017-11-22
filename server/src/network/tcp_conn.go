package network

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	ConnInit     = 0
	ConnRegister = 1
	ConnClose    = 2

	ConnPingMaxTime = 15
	ConnRegisterMaxTime = 30

	IsEncrypt = true
	IsZip     = true

	ReadClose  = 1
	WriteClose = 2
	ForceClose = 3
)

var AES_KEY = []byte("a4Y5T7G901A3TOZe")

type TagType int8

const (
	RawTag = TagType(0)
	ZipTag = TagType(1) << iota
	EncryptTag
)

type TcpConn struct {
	Type           int
	ID             int
	LastPing       int64
	initTime       int64
	waitGroup      sync.WaitGroup
	closeLock      sync.Once
	onCloseLock    sync.Once
	Conn           *net.TCPConn // the raw connection
	Status         int
	AttachID       int
	ExtraData      interface{}
	reader         *bufio.Reader
	handler        Handler
	writeChan      chan []byte
	closeWaitGroup sync.WaitGroup
	aes            AesEncrypt
}

//进行zlib压缩
func zlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func zlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

type Handler interface {
	OnMessage(conn *TcpConn, bytes []byte)
	OnClose(conn *TcpConn)
}

func NewTcpConn(conn *net.TCPConn, id int, handler Handler) *TcpConn {
	tcpConn := &TcpConn{
		ID:       id,
		handler:  handler,
		Conn:     conn,
		LastPing: time.Now().Unix(),
	}
	tcpConn.Init()
	go tcpConn.ReadLoop()
	go tcpConn.WriteLoop()
	return tcpConn
}

func Connect(id int, addr string, handler Handler) (*TcpConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("ConnectServer fail! add:", addr, " err:", err)
		return nil, err
	}
	return NewTcpConn(conn.(*net.TCPConn), id, handler), nil
}

func (c *TcpConn) Init() {
	c.reader = bufio.NewReader(c.Conn)
	c.writeChan = make(chan []byte, WRITE_CHAN_SIZE)
	c.Status = ConnInit
	c.aes.SetKey(AES_KEY)
	c.initTime = time.Now().Unix()
}

func (c *TcpConn) GetID() int {
	return c.ID
}

func (c *TcpConn) OnClose(closeType int) {
	c.onCloseLock.Do(func() {
		c.Status = ConnClose
		close(c.writeChan)
		c.Conn.Close()
		if closeType == ForceClose {
			c.WaitClose()
		}
		c.handler.OnClose(c)
	})
}

func (c *TcpConn) Close() {
	c.OnClose(ForceClose)
}

func (c *TcpConn) WaitClose() {
	c.closeWaitGroup.Wait()
}

//                  |1byte|4byte   | length byte
//  message format  |tag  | length | data
//

func (c *TcpConn) ReadLoop() {
	// read length
	var size uint32
	var sizeBytes = make([]byte, HEADER_SIZE)
	c.closeWaitGroup.Add(1)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		log.Println("ReadLoop close")
		c.closeWaitGroup.Done()
		c.OnClose(ReadClose)
	}()

	for {
		if _, err := io.ReadFull(c.reader, sizeBytes); err != nil {
			log.Println(err)
			break
		}
		if size = binary.BigEndian.Uint32(sizeBytes[1:]); size > MAX_MESSAGE_SIZE {
			log.Println("error package size:", size)
			break
		}

		buff := make([]byte, size)
		if _, err := io.ReadFull(c.reader, buff[:]); err != nil {
			log.Println(err)
			break
		}

		data := buff[:]
		tag := TagType(sizeBytes[0])
		if tag != RawTag {
			if (tag & EncryptTag) != 0 {
				var err error
				data, err = c.aes.Decrypt(data)
				if err != nil {
					log.Println(err)
					break
				}
			}
			if (tag & ZipTag) != 0 {
				data = zlibUnCompress(data)
			}
		}
		c.handler.OnMessage(c, data)
		c.LastPing = time.Now().Unix()
	}
}

func (c *TcpConn) WriteLoop() {
	c.closeWaitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		log.Println("WriteLoop close")
		c.closeWaitGroup.Done()
		c.OnClose(WriteClose)
	}()
	sec := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sec.C: //ping check
			if time.Now().Unix() - c.LastPing > ConnPingMaxTime {
				log.Println("ping check time out!")
				return
			}
			if (c.Status == ConnInit) && (time.Now().Unix() - c.initTime > ConnRegisterMaxTime) {
				log.Println("register check time out!")
				return
			}
		case bytes, ok := <-c.writeChan:
			if !ok {
				return
			}

			sendBytes := bytes[:]
			tag := RawTag
			if IsZip {
				zipBytes := zlibCompress(bytes)
				if len(zipBytes) < len(bytes) {
					log.Println("zip len:", len(zipBytes), " src len:", len(bytes))
					tag += ZipTag
					sendBytes = zipBytes[:]
				}
			}

			if IsEncrypt {
				tag += EncryptTag
				var err error
				sendBytes, err = c.aes.Encrypt(sendBytes)
				if err != nil {
					log.Println(err)
					break
				}
			}

			buff := make([]byte, HEADER_SIZE+len(sendBytes))
			buff[0] = byte(tag)
			binary.BigEndian.PutUint32(buff[1:HEADER_SIZE], uint32(len(sendBytes)))
			copy(buff[HEADER_SIZE:], sendBytes)

			err := c.OnWrite(buff)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *TcpConn) SendPacket(bytes []byte) {
	if c.Status == ConnClose {
		return
	}
	c.writeChan <- bytes
}

func (c *TcpConn) OnWrite(bytes []byte) error {
	sendSize, err := c.Conn.Write(bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	if sendSize != len(bytes) {
		panic("error send len!")
	}
	return nil
}
