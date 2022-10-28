package accesskey

// AccessKey represents the AccessKey entity
type AccessKey string

func (a AccessKey) String() string {
	return string(a)
}
