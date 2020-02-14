package retrier

// LastErrorWrapper always return the last error received, unmodified. If
// the last error was nil, it will return nil as well
type LastErrorWrapper struct {
}

// NewLastErrorWrapper returns an instance of LastErrorWrapper
func NewLastErrorWrapper() *LastErrorWrapper {
	return &LastErrorWrapper{}
}

// WrapError will always return @err unmodified. If @err is nil, returns nil.
func (w *LastErrorWrapper) WrapError(_ int, err error) error {
	return err
}
