package main

import (
	"./client"
	"log"
	"runtime"
)

var closeChan chan bool
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    log.SetFlags(log.Lshortfile|log.Ltime)
    log.Println("client begin start!")
    client.Start()
    <-closeChan
}