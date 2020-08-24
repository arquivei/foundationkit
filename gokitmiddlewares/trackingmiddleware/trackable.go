package trackingmiddleware

import (
	"github.com/arquivei/foundationkit/trace"
)

// Traceable is an interface of something that has trace.
//
// This is intentend to be used in places that you don't have a transport middleware
// to initialize the trace but you have the trace in the request, so you can make the
// request Traceable and have this middleware extracting the Trace and placing on the
// cotnext.
type Traceable interface {
	Trace() trace.Trace
}
