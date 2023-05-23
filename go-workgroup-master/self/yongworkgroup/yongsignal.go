package yongworkgroup

import (
	"os"
	"os/signal"
)

func YongSignal(sig ...os.Signal) RunYongFunc {
	return func(stop <-chan struct{}) error {
		if len(sig) == 0 {
			sig = append(sig, os.Interrupt)
		}
		done := make(chan os.Signal, len(sig))

		defer close(done)

		signal.Notify(done, sig...)

		defer signal.Stop(done)

		select {
		case <-stop:
		case <-done:
		}

		return nil
	}
}
