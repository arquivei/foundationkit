package stakeholder

import (
	"encoding/json"

	"github.com/arquivei/foundationkit/errors"
)

// Stakeholder is used to represent a person or company involved
// in the resources of the system, interested in it's results.
type Stakeholder string

func (s Stakeholder) String() string {
	return string(s)
}

//MarshalJSON serializes the stakeholder value as a JSON value
func (s Stakeholder) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

//UnmarshalJSON deserializes a JSON value into a Stakeholder value
func (s *Stakeholder) UnmarshalJSON(b []byte) error {
	const op = errors.Op("stakeholder.UnmarshalJSON")
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return errors.E(op, err)
	}
	*s, err = Parse(v)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

// GetType returns the stakeholder type
func GetType(s Stakeholder) Type {
	return getType(string(s))
}

func getType(s string) Type {
	switch len(s) {
	case 14:
		return TypeCompany
	case 11:
		return TypePerson
	default:
		return TypeUnknown
	}
}

// Parse attempts to create a stakeholder from it's ID
func Parse(s string) (Stakeholder, error) {
	const op = errors.Op("stakeholder.Parse")

	t := getType(s)
	switch t {
	case TypeCompany:
		return NewCNPJ(s)
	case TypePerson:
		return NewCPF(s)
	default:
		return "", errors.E(op, newInvalidStakeholderTypeError(t))
	}
}

// NewCPF creates a new stakeholder with TypePerson
func NewCPF(cpf string) (Stakeholder, error) {
	if err := CheckCPF(cpf); err != nil {
		return "", err
	}
	return Stakeholder(cpf), nil
}

// NewCNPJ creates a new stakeholder with TypeCompany
func NewCNPJ(cnpj string) (Stakeholder, error) {
	if err := CheckCNPJ(cnpj); err != nil {
		return "", err
	}
	return Stakeholder(cnpj), nil
}

// GetCPFCNPJ returns the CPF or CNPJ of the given stakeholder.
// Since a stakeholder cannot have both, one of them will be returned
// empty, or both if stakeholder is of unknown type
func GetCPFCNPJ(s Stakeholder) (cpf, cnpj string) {
	switch GetType(s) {
	case TypePerson:
		cpf = string(s)
	case TypeCompany:
		cnpj = string(s)
	}

	return
}
