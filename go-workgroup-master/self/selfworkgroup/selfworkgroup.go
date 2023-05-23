package selfworkgroup

type RunSelfFunc func(<-chan struct{}) error

type SelfGroup struct {
	fnself []RunSelfFunc
}

func (sg *SelfGroup) Add(sfn RunSelfFunc) {
	sg.fnself = append(sg.fnself, sfn)
}

func (sg *SelfGroup) Run() error {
	if len(sg.fnself) == 0 {
		return nil
	}

	stop := make(chan struct{})
	done := make(chan error, len(sg.fnself))

	defer close(done)

	for _, fnself := range sg.fnself {
		go func(fnself RunSelfFunc) {
			done <- fnself(stop)
		}(fnself)
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
