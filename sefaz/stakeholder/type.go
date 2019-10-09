package stakeholder

import "strconv"

// Type is a strong type to classify a Stakeholder
type Type int

const (
	// NOTE : DO NOT REORDER OR DELETE VALUES
	// Type values are being used on serialization and deserialization of
	// messages

	// TypeUnknown is the default value, and usually means an error
	TypeUnknown Type = iota
	// TypePerson is used to represent 'Pessoa Fisica'
	TypePerson
	// TypeCompany is used to represent 'Pessoa Juridica'
	TypeCompany
)

// TypeText returns the stakeholder type as a text
func TypeText(t Type) string {
	switch t {
	case TypePerson:
		return "cpf"
	case TypeCompany:
		return "cnpj"
	case TypeUnknown:
		return "unkown"
	}
	return "unexpected: " + strconv.Itoa(int(t))
}
