package consul
import(
	"sync"
)

type CloseWait struct {
	waitGroup sync.WaitGroup
}
func (c *CloseWait) Add() {
	c.waitGroup.Add(1)
}

func (c *CloseWait) Done() {
	c.waitGroup.Done()
}

func (c *CloseWait) Wait() {
	c.waitGroup.Wait()
}