package nsu

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/arquivei/foundationkit/errors"
)

// NSU is "Número Sequencial Único" and it's an offset used by sefaz to manage the NFes
type NSU string

func (nsu NSU) String() string {
	return fmt.Sprintf("%015s", string(nsu))
}

//MarshalJSON serializes the NSU value as a JSON value
func (nsu NSU) MarshalJSON() ([]byte, error) {
	// I don't capture the error here and wraps it with op because it's impossible (probably)
	// that we can reproduce said error in a unit test and we would have some untestable code.
	return json.Marshal(nsu.String())
	// Just for reference, here is the code I previously wrote:
	/*
		const op = errors.Op("nsu.MarshalJSON")
		b, err := json.Marshal(nsu.String())
		if err != nil {
			return nil, errors.E(op, err)
		}
		return b, nil
	*/
}

//UnmarshalJSON deserializes a JSON value into a NSU value
func (nsu *NSU) UnmarshalJSON(b []byte) error {
	const op = errors.Op("nsu.UnmarshalJSON")
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return errors.E(op, err)
	}
	*nsu, err = Parse(s)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// Zero is a zero valued NSU.
const Zero = NSU("000000000000000")

// Parse instantiates a new nsu from @nsu string
func Parse(nsu string) (NSU, error) {
	const op = errors.Op("nsu.Parse")

	if len(nsu) == 0 {
		return "", errors.E(op, "nsu is empty")
	}
	if len(nsu) > 15 {
		return "", errors.E(op, "nsu has more than 15 digits")
	}

	for i := 0; i < len(nsu); i++ {
		if nsu[i]-'0' > 9 {
			return "", errors.E(op, ErrCannotParse)
		}
	}

	return NSU(fmt.Sprintf("%015s", nsu)), nil
}

// MustParse calls Parse function and panics on error
func MustParse(s string) NSU {
	const op = errors.Op("nsu.MustParse")

	nsu, err := Parse(s)
	if err != nil {
		panic(errors.E(op, err))
	}
	return nsu
}

// AsInt converts a NSU into an Integer. This function panics if the NSU is not an integer
func AsInt(nsu NSU) int {
	const op = errors.Op("nsu.AsInt")
	i, err := strconv.Atoi(string(nsu))
	if err != nil {
		panic(err)
	}
	return i
}

// Compare two NSU's by using this function. NSU's will be compared after
// being added the padding 0'es. Returns 0 if @source is the same as @compare,
// 1 if @source is bigger and -1 if @source is smaller
func Compare(source, compare NSU) int {
	return strings.Compare(source.String(), compare.String())
}
