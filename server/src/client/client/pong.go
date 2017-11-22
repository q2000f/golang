package client
import(
    "../../protocol"
    "../../network"
    "encoding/json"
    "log"
)

func Pong(conn *network.TcpConn, body []byte) {
    ping := protocol.Ping{}
    json.Unmarshal(body, &ping)
    log.Println("Pong:", string(body))
}