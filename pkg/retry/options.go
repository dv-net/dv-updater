package retry

import (
	"context"
	"time"
)

type Option func(*Retry)

func WithMaxAttempts(maxAttempts int) Option {
	return func(r *Retry) {
		r.maxAttempts = maxAttempts
	}
}

func WithPolicy(policy Policy) Option {
	return func(r *Retry) {
		r.policy = policy
	}
}

func WithDelay(delay time.Duration) Option {
	return func(r *Retry) {
		r.delay = delay
	}
}

func WithDebug(debug bool) Option {
	return func(r *Retry) {
		r.debug = debug
	}
}

func WithContext(ctx context.Context) Option {
	return func(r *Retry) {
		r.ctx = ctx
	}
}

func (r *Retry) SetMaxAttempts(maxAttempts int) *Retry {
	r.maxAttempts = maxAttempts
	return r
}

func (r *Retry) SetPolicy(policy Policy) *Retry {
	r.policy = policy
	return r
}

func (r *Retry) SetDelay(delay time.Duration) *Retry {
	r.delay = delay
	return r
}

func (r *Retry) SetDebug(debug bool) *Retry {
	r.debug = debug
	return r
}
