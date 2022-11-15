package go_routine

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

type counterRace struct {
	sync.Mutex
	val int
}

func (c *counterRace) Add(x int) {
	c.Lock()
	c.val++
	c.Unlock()
}

func (c *counterRace) Value() (x int) {
	return c.val
}

func TestRace(t *testing.T) {
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	var meter counterRace

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			for j := 0; j < 1000; j++ {
				meter.Add(1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(meter.Value())
}
