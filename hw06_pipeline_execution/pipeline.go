package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		out := make(Bi)
		close(out)
		return out
	}

	for _, stage := range stages {
		in = runStage(stage(in), done)
	}

	return in
}

func runStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range in {
				continue
			}
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}
