package cuf

import (
	"encoding/json"
	"github.com/arquivei/foundationkit/errors"
)

// CUF stands for "Codigo da unidade federativa" and it
// is strongly associated with Stakeholder.
type CUF string

// New validate and returns either (CUF, nil) if the given cUF
// is valid or (empty CUF, error).
func New(uf string) (CUF, error) {
	cUF := CUF(uf)
	if cUF == "" {
		return "", errors.New("missing cUF")
	}
	if isValidUF(uf) {
		return cUF, nil
	}
	return "", errors.New("invalid cUF")
}

func (c CUF) String() string {
	return string(c)
}

//MarshalJSON serializes the CUF value as a JSON value
func (s CUF) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

//UnmarshalJSON deserialize a JSON value into a CUF value
func (s *CUF) UnmarshalJSON(b []byte) error {
	const op = errors.Op("CUF.UnmarshalJSON")
	var v string

	err := json.Unmarshal(b, &v)
	if err != nil {
		return errors.E(op, err)
	}
	if *s, err = New(v); err != nil {
		return err
	}

	return nil
}

func isValidUF(cUF string) bool {
	if len(cUF) != 2 {
		return false
	}
	switch cUF[0] {
	case '1':
		switch cUF[1] {
		case '0', '8', '9':
			return false
		default:
			return true
		}
	case '2':
		switch cUF[1] {
		case '0':
			return false
		default:
			return true
		}
	case '3':
		switch cUF[1] {
		case '1', '2', '3', '5':
			return true
		default:
			return false
		}
	case '4':
		switch cUF[1] {
		case '1', '2', '3':
			return true
		default:
			return false
		}
	case '5':
		switch cUF[1] {
		case '0', '1', '2', '3':
			return true
		default:
			return false
		}
	default:
		return false
	}
}
