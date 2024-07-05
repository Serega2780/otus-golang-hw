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

type Sync struct {
	m  map[int]interface{}
	mu sync.Mutex
}

type Mark struct {
	pos int
	v   interface{}
}

func NewSync() *Sync {
	s := new(Sync)
	s.m = make(map[int]interface{})
	return s
}

var GrCount = 8

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.

	preparedC := make(Bi)
	doneC := make(Bi)
	grCount := 0

	s := NewSync()

	prepare(in, preparedC)

	runWorkers(s, preparedC, stages, doneC)

	for {
		select {
		case <-done:
			return calcResult(s)
		case <-doneC:
			grCount++
			if grCount == GrCount {
				return calcResult(s)
			}
		}
	}
}

func calcResult(s *Sync) Bi {
	s.mu.Lock()
	result := make(Bi, len(s.m))
	if len(s.m) > 0 {
		for i := 0; i < len(s.m); i++ {
			v := s.m[i+1]
			if v != nil {
				result <- s.m[i+1]
			}
		}
	}
	s.mu.Unlock()
	close(result)
	return result
}

func runWorkers(s *Sync, preparedC Bi, stages []Stage, doneC Bi) {
	for i := 0; i < GrCount; i++ {
		go worker(s, preparedC, stages, doneC)
	}
}

func prepare(in In, preparedC Bi) {
	go func() {
		i := 0
		for v := range in {
			i++
			preparedC <- Mark{i, v}
		}
		close(preparedC)
	}()
}

func worker(s *Sync, preparedC In, stages []Stage, doneC Bi) {
	for v := range preparedC {
		key := v.(Mark).pos
		tmp := steps(stages, v.(Mark).v)
		s.mu.Lock()
		s.m[key] = tmp
		s.mu.Unlock()
	}
	doneC <- struct{}{}
}

func steps(stages []Stage, v interface{}) interface{} {
	tmpV := v
	for i := 0; i < len(stages); i++ {
		tmpV = step(stages[i], tmpV)
	}
	return tmpV
}

func step(stage Stage, v interface{}) interface{} {
	tmpC := make(Bi, 1)
	tmpC <- v
	close(tmpC)
	tmp := <-stage(tmpC)
	return tmp
}
