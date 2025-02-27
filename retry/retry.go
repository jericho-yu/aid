package retry

import (
	"context"
	"math/rand"
	"time"
)

// Do 重试机制
func Do(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Do(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

// WithContext 带上下文的重试机制
func WithContext(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			select {
			case <-time.After(sleep):
				return WithContext(ctx, attempts, 2*sleep, fn) // 指数退避
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	return nil
}

// WithContextAndJitter 带上下文和随机退避的重试机制
func WithContextAndJitter(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			// 加入随机退避
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter

			select {
			case <-time.After(sleep):
				return WithContextAndJitter(ctx, attempts, 2*sleep, fn) // 指数退避
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	return nil
}
