package ping

import "time"

// Request Request
type Request struct {
	Num   int
	Sleep time.Duration
}
