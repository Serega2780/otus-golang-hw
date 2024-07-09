package hw06pipelineexecution

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
			defer close(out)
			for {
				select {
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				case <-done:
					return
				}
			}
		}(in, out)
		in = stages[i](out)
	}
	result = in
	return result
}
