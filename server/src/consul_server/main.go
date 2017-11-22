package main
import (
    "./consul"
    "runtime"
    "log"
)
var closeChan chan bool
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.SetFlags(log.Lshortfile|log.Ltime)
    log.Println("consul server begin start!")
    consul.Start()
    <-closeChan
}


