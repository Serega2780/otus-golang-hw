package hw06pipelineexecution

import (
	"fmt"
	"sync/atomic"
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	var stop atomic.Bool
	start := make(Bi)
	go func() {
		start <- struct{}{}
		_, ok := <-done
		if !ok {
			t := time.Now()
			fmt.Printf("%t\n", ok)
			fmt.Println(t.Format("15:04:05.000"))
			stop.CompareAndSwap(false, true)
		}
	}()
	<-start
	return f(stages, in, &stop)
}

func f(stages []Stage, in In, stop *atomic.Bool) Out {
	for i := 0; i < len(stages); i++ {
		t := time.Now()
		fmt.Println(t.Format("15:04:05.000"))
		fmt.Printf("stop %t\n", stop.Load())
		if stop.Load() {
			return in
		}
		in = stages[i](in)
	}
	return in
}
