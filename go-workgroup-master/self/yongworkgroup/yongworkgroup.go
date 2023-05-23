package yongworkgroup

type RunYongFunc func(<-chan struct{}) error

type YongGroup struct {
	fns []RunYongFunc
}

func (g *YongGroup) Add(fn RunYongFunc) {
	g.fns = append(g.fns, fn)
}

func (g *YongGroup) RunYong() error {
	if len(g.fns) == 0 {
		return nil
	}

	stop := make(chan struct{})
	done := make(chan error, len(g.fns))

	defer close(done)

	for _, fn := range g.fns {
		go func(fn RunYongFunc) {
			done <- fn(stop)
		}(fn)
	}

	var err error
	for i := 0; i < cap(done); i++ {
		if err == nil {
			err = <-done
		} else {
			<-done
		}
		if i == 0 {
			close(stop)
		}
	}
	return err
}
