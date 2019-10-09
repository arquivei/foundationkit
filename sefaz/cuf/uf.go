package cuf

import "errors"

// Cuf stands for "Codigo da unidade federativa" and it
// is strongly associated with Stakeholder.
type Cuf string

// New validate and returns either (Cuf, nil) if the given cUF
// is valid or (empty Cuf, error).
func New(cUF string) (Cuf, error) {
	if cUF == "" {
		return Cuf(""), errors.New("missing cUF")
	}
	if isValidUF(cUF) {
		return Cuf(cUF), nil
	}
	return Cuf(""), errors.New("invalid cUF")
}

func (c Cuf) String() string {
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
