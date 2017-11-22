package consul
var timerInstance *Timer
func Start() {
	timerInstance = NewTimer()
	InitServerMgr()
	InitNetWork()
}
