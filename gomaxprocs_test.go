package go_routine

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetGomaxprocs(t *testing.T) {
	totalCpu := runtime.NumCPU()
	fmt.Println("Total CPU: ", totalCpu)

	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println("Total Threads: ", totalThread)
}
