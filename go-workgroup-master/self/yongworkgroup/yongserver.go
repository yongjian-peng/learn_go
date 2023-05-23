package yongworkgroup

func YongServer(serve func() error, shutdown func() error) RunYongFunc {
	return func(stop <-chan struct{}) error {
		done := make(chan error)

		defer close(done)

		go func() {
			done <- serve()
		}()

		select {
		case err := <-done:
			return err
		case <-stop:
			err := shutdown()
			if err == nil {
				err = <-done
			} else {
				<-done
			}
			return err
		}
	}
}
