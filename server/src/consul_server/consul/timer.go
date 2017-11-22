package consul

import "time"

//type IDayTimer interface {
//    OnPerDay()
//}
//
//type IHourTimer interface {
//    OnPerHour()
//}
//
//type IMinuteTimer interface {
//    OnPerMinute()
//}

type ISecondTimer interface {
    OnPerSecond()
}

type Timer struct {
    CloseWait
    closeChan chan bool
    SecondTimerArr []ISecondTimer
}

func NewTimer() *Timer {
    timer := Timer{}
    timer.Init()
    return &timer
}

func (c *Timer) Init() {
    c.SecondTimerArr = []ISecondTimer{}
    go c.Loop()
}

func (c *Timer) Close() {
    c.closeChan  <- true
}

func (c *Timer) InsertSecondTimer(timer ISecondTimer) {
    c.SecondTimerArr = append(c.SecondTimerArr, timer)
}

func (c *Timer) Loop() {
    c.CloseWait.Add()
    second := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-c.closeChan:
            break
        case <-second.C:
            for _, timer := range c.SecondTimerArr {
                timer.OnPerSecond()
            }
        }
    }
    c.CloseWait.Done()
}