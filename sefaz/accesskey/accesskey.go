package accesskey

// Key represents the AccessKey entity
type Key string

// AccesskeyValidator is the interface that performs accesskey validation
type Validator interface {
	Check(Key) error
}

type validator struct{}

// NewValidator returns a real validator
func NewValidator() Validator {
	return &validator{}
}

func (a Key) String() string {
	return string(a)
}
