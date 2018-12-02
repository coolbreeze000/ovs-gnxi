package retry

import (
	"time"
)

// Function signature of retry if function
type RetryIfFunc func(error) bool

// Function signature of OnRetry function
// n = count of attempts
type OnRetryFunc func(n uint, err error)

type config struct {
	attempts uint
	delay    time.Duration
	units    time.Duration
	onRetry  OnRetryFunc
	retryIf  RetryIfFunc
}

// Option represents an option for retry.
type Option func(*config)

// Attempts set count of retry
// default is 10
func Attempts(attempts uint) Option {
	return func(c *config) {
		c.attempts = attempts
	}
}

// Delay set delay between retry
// default are 1e5 units
func Delay(delay time.Duration) Option {
	return func(c *config) {
		c.delay = delay
	}
}

// Units set unit of delay (probably only for tests purpose)
// default are microsecond
func Units(units time.Duration) Option {
	return func(c *config) {
		c.units = units
	}
}

// OnRetry function callback are called each retry
//
// log each retry example:
//
//	retry.Do(
//		func() error {
//			return errors.New("some error")
//		},
//		retry.OnRetry(func(n unit, err error) {
//			log.Printf("#%d: %s\n", n, err)
//		}),
//	)
func OnRetry(onRetry OnRetryFunc) Option {
	return func(c *config) {
		c.onRetry = onRetry
	}
}

// RetryIf controls whether a retry should be attempted after an error
// (assuming there are any retry attempts remaining)
//
// skip retry if special error example:
//
//	retry.Do(
//		func() error {
//			return errors.New("special error")
//		},
//		retry.RetryIf(func(err error) bool {
//			if err.Error() == "special error" {
//				return false
//			}
//			return true
//		})
//	)
func RetryIf(retryIf RetryIfFunc) Option {
	return func(c *config) {
		c.retryIf = retryIf
	}
}
