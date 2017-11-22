package protocol

import "encoding/binary"
const(
    HEADER_SIZE = 4
)

//Packet is pack message
func Packet(bytes []byte)[]byte {
    var sendBytes = make([]byte, HEADER_SIZE + len(bytes))
    binary.BigEndian.PutUint32(sendBytes[0:HEADER_SIZE], uint32(len(bytes)))
    copy(sendBytes[HEADER_SIZE:], bytes)
    return sendBytes
}
