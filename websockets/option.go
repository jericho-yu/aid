package websockets

import "time"

type (
	Heart struct {
		timer *time.Ticker
	}

	Timeout struct {
		interval time.Duration
		timer    string
	}
)

func NewHeart(interval time.Duration) {
	time.NewTicker(interval)
}
