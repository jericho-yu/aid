package retry

import (
	"context"
	"math/rand"
	"time"
)

// Retry 重试机制
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

// RetryWithContext 带上下文的重试机制
func RetryWithContext(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			select {
			case <-time.After(sleep):
				return RetryWithContext(ctx, attempts, 2*sleep, fn) // 指数退避
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	return nil
}

// RetryWithContextAndJitter 带上下文和随机退避的重试机制
func RetryWithContextAndJitter(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			// 加入随机退避
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter

			select {
			case <-time.After(sleep):
				return RetryWithContextAndJitter(ctx, attempts, 2*sleep, fn) // 指数退避
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	return nil
}
