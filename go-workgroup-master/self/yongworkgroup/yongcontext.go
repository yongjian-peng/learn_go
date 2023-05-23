package yongworkgroup

import "context"

func YongContext(ctx context.Context) RunYongFunc {
	return func(stop <-chan struct{}) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stop:
			return nil
		}
	}
}
