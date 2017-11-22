package server
import(
    "../../protocol"
    "../../network"
    "log"
    "time"
)

func Ping(conn *network.TcpConn, body []byte) (interface{}, error) {
    req := protocol.Ping{}
    err := protocol.Unmarshal(&req, body)
    res := &protocol.Pong{}
    if err != nil {
        result := []byte("error!");
        return result, nil
    }
    res.ClientTime = req.ClientTime
    res.ServerTime = time.Now().UnixNano()
    log.Println("ping:", string(body))
    return res, nil
}

