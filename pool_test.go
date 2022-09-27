package go_routine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//func RunAsyncPool(group *sync.WaitGroup, pool *sync.Pool) {
//	defer group.Done()
//
//	group.Add(1)
//	data := pool.Get()
//	fmt.Println(data)
//	//time.Sleep(1 * time.Second)
//	pool.Put(data)
//	//fmt.Println("Hello")
//	time.Sleep(11 * time.Second)
//}

func TestPool(t *testing.T) {
	group := &sync.WaitGroup{}
	pool := &sync.Pool{
		New: func() interface{} {
			return "New"
		},
	}

	pool.Put("Raka")
	pool.Put("Janitra")

	for i := 0; i < 10; i++ {
		go func() {
			group.Add(1)
			data := pool.Get()
			fmt.Println(data)
			time.Sleep(1 * time.Second)
			pool.Put(data)
		}()
	}
	//group.Wait()
	time.Sleep(11 * time.Second)
	fmt.Println("Selesai")
}
