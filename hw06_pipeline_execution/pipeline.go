package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	var wg sync.WaitGroup
	result := make([]interface{}, 0)

	wg.Add(1)
	go worker(&wg, stages, in, done, &result)
	wg.Wait()

	return calcResult(result)
}

func calcResult(result []interface{}) Bi {
	out := make(Bi, len(result))
	defer close(out)
	for _, r := range result {
		out <- r
	}
	return out
}

func worker(wg *sync.WaitGroup, stages []Stage, in In, done In, result *[]interface{}) {
	defer wg.Done()
	for i := 0; i < len(stages); i++ {
		in = stages[i](in)
	}
	for {
		select {
		case <-done:
			return
		case tmp, ok := <-in:
			if !ok {
				return
			}
			*result = append(*result, tmp)
		}
	}
}
