package cuf

import "errors"

// CUF stands for "Codigo da unidade federativa" and it
// is strongly associated with Stakeholder.
type CUF string

// New validate and returns either (CUF, nil) if the given cUF
// is valid or (empty CUF, error).
func New(cUF string) (CUF, error) {
	if cUF == "" {
		return CUF(""), errors.New("missing cUF")
	}
	if isValidUF(cUF) {
		return CUF(cUF), nil
	}
	return CUF(""), errors.New("invalid cUF")
}

func (c CUF) String() string {
	return string(c)
}

func isValidUF(uf string) bool {
	if len(uf) != 2 {
		return false
	}
	switch uf[0] {
	case '1':
		switch uf[1] {
		case '0', '8', '9':
			return false
		default:
			return true
		}
	case '2':
		switch uf[1] {
		case '0':
			return false
		default:
			return true
		}
	case '3':
		switch uf[1] {
		case '1', '2', '3', '5':
			return true
		default:
			return false
		}
	case '4':
		switch uf[1] {
		case '1', '2', '3':
			return true
		default:
			return false
		}
	case '5':
		switch uf[1] {
		case '0', '1', '2', '3':
			return true
		default:
			return false
		}
	default:
		return false
	}
}
