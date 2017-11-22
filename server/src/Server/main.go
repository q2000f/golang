package main

import (
	"log"
	"runtime"
	"os/signal"
	"os"
	"./server"
)
var closeChan chan bool

func WaitForSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)
	<-c
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("server begin start!")
	server.Start();
	WaitForSignal()
	server.Stop();
}
