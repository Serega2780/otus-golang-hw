package hw06pipelineexecution

import (
	"sync/atomic"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	var result Out
	for i := 0; i < len(stages); i++ {
		out := make(Bi)
		go func(in In, out Bi) {
			var isOutClosed atomic.Bool
			isOutClosed.Store(false)
			for {
				select {
				case v, ok := <-in:
					if !ok {
						if !isOutClosed.Load() {
							close(out)
						}
						return
					}
					if !isOutClosed.Load() {
						out <- v
					}
				case <-done:
					if !isOutClosed.Swap(true) {
						close(out)
					}
				}
			}
		}(in, out)
		in = stages[i](out)
	}
	result = in
	return result
}
