package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {}
		}()

		if in == nil || len(stages) == 0 {
			return
		}

		for _, s := range stages {
			select {
			case <-done:
				return
			default:
				in = s(in)
			}
		}

		for {
			select {
			case <-done:
				return
			default:
				select{
				case <-done:
					return
				case v, ok := <- in:
					if !ok {
						return
					}
					out <- v
				}
			}
		}
	}()
	return out
}