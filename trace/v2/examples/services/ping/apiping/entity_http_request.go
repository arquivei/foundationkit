package apiping

import "time"

// Request Request
type Request struct {
	Num   int           `json:"num"`
	Sleep time.Duration `json:"sleep"`
}
