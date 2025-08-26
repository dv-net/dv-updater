package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Retry struct {
	ctx         context.Context
	maxAttempts int
	policy      Policy
	delay       time.Duration
	debug       bool
}

var (
	ErrRetry = errors.New("retry")
	ErrExit  = errors.New("exit")
)

func New(opts ...Option) *Retry {
	r := &Retry{
		maxAttempts: 5,
		policy:      PolicyBackoff,
		delay:       1 * time.Second,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Retry) Do(fn func() error) error {
	switch r.policy {
	case PolicyLinear:
		return r.linearRetry(fn)
	case PolicyBackoff:
		return r.backoffRetry(fn)
	case PolicyInfinite:
		return r.infiniteRetry(fn)
	default:
		return fmt.Errorf("unsupported retry policy")
	}
}

func (r *Retry) linearRetry(fn func() error) error {
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrExit) {
			return err
		}

		if attempt < r.maxAttempts {
			if r.debug {
				fmt.Printf("linear Retry attempt %d failed, retrying in %s...\n", attempt, r.delay)
			}
			time.Sleep(r.delay)
		}
	}
	return fmt.Errorf("linear retry failed after %d attempts", r.maxAttempts)
}

func (r *Retry) backoffRetry(fn func() error) error {
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrExit) {
			return err
		}

		if attempt < r.maxAttempts {
			delay := r.delay * (1 << (attempt - 1)) // Increase *2 every attempt
			if r.debug {
				fmt.Printf("backoff Retry attempt %d failed, retrying in %s...\n", attempt, delay)
			}
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("backoff retry failed after %d attempts", r.maxAttempts)
}

func (r *Retry) infiniteRetry(fn func() error) error {
	if r.ctx == nil {
		return fmt.Errorf("infinite retry cannot be initialized without ctx")
	}

	resCh := make(chan error, 1)
	go func() {
		defer close(resCh)
		for {
			select {
			case <-r.ctx.Done():
				return
			default:
				err := fn()
				if err == nil {
					return
				}

				if errors.Is(err, ErrExit) {
					resCh <- err
					return
				}

				if r.debug {
					fmt.Printf("initnite retry attempt\n")
				}
				time.Sleep(r.delay)
			}
		}
	}()

	return <-resCh
}
