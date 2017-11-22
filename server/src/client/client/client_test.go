package client
import (
    "testing"
    "runtime"
    "../../protocol"
)

func BenchmarkConnect(b *testing.B) {
    runtime.GOMAXPROCS(runtime.NumCPU())
    client := NewClient("localhost:1236")
    reg := protocol.IRegisterServer{}
    reg.ServerID = 1001
    reg.ServerType = 1
    reg.ServerName = "game"
    client.RegisterServer(reg)
}