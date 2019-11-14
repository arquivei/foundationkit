package cuf

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/arquivei/foundationkit/errors"
)

// CUF stands for "Código da Unidade Federativa"
type CUF struct {
	initialized bool
	value       uint8
}

// New validate and returns either (CUF, nil) if the given cUF
// is valid or (empty CUF, error).
func New(uf string) (CUF, error) {
	const op = errors.Op("cuf.New")
	ufInt, err := parseUF(uf)

	if err != nil {
		return CUF{}, errors.E(op, err)
	}
	return CUF{true, ufInt}, nil
}

// MustNew returns CUF if the given cUF is valid or panic
func MustNew(uf string) CUF {
	cuf, err := New(uf)
	if err != nil {
		panic(err)
	}
	return cuf
}

func (c CUF) String() string {
	if !c.initialized {
		return ""
	}
	return strconv.Itoa(int(c.value))
}

//MarshalJSON serializes the CUF value as a JSON value
func (c CUF) MarshalJSON() ([]byte, error) {
	const op = errors.Op("CUF.MarshalJSON")

	if !c.initialized {
		return nil, errors.E(op, "CUF not initialized")
	}
	return json.Marshal(c.String())
}

//UnmarshalJSON deserialize a JSON value into a CUF value
func (c *CUF) UnmarshalJSON(b []byte) error {
	const op = errors.Op("CUF.UnmarshalJSON")
	var v string

	err := json.Unmarshal(b, &v)
	if err != nil {
		return errors.E(op, err)
	}
	if *c, err = New(v); err != nil {
		return errors.E(op, err)
	}

	return nil
}

func IsValid(cUF CUF) bool {
	return cUF.initialized
}

func parseUF(cUF string) (uint8, error) {
	if len(cUF) != 2 {
		return 0, errors.Errorf("input cUF should have 2 digits: %s", cUF)
	}

	ufInt64, err := strconv.ParseUint(cUF, 10, 8)
	if err != nil {
		return 0, errors.Errorf("cUF could not be converted to integer: %s",
			cUF)
	}

	switch ufInt64 {
	case 11, 12, 13, 14, 15, 16, 17, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31, 32,
		33, 35, 41, 42, 43, 50, 51, 52, 53:
		return uint8(ufInt64), nil
	default:
		return 0, errors.Errorf(fmt.Sprintf("invalid cUF code: %s", cUF))
	}
}
