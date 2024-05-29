package accesskey

import (
	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/sefaz/stakeholder"
)

// Check checks if an access key is valid. It returns an error with an specific code based on
// the validation problem
func Check(accessKey AccessKey) error {
	const op errors.Op = "accesskey.Check"

	err := validate(accessKey)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

// CheckNFF checks if an access key is valid with more restrictive rules also considering NFF.
// It returns an error with an specific code based on the validation proble
func CheckNFF(accessKey AccessKey) error {
	const op errors.Op = "accesskey.CheckWithNFF"

	err := validate(accessKey)
	if err != nil {
		return errors.E(err, op)
	}

	if isNFF(accessKey) {
		err := validateNFF(accessKey)
		if err != nil {
			return errors.E(err, op)
		}
	}

	return nil
}

// validate execute all sub-routines necessary to perform a full accesskey validation for regular NFes
func validate(accessKey AccessKey) error {
	const op errors.Op = "validate"

	if accessKey == "" {
		return errors.E(ErrEmptyAccessKey, ErrCodeEmptyAccessKey, op)
	}

	if len(accessKey) != 44 {
		return errors.E(ErrInvalidLength, ErrCodeInvalidLength, op)
	}

	if !isDigitOnly(accessKey) {
		return errors.E(ErrInvalidCharacter, ErrCodeInvalidCharacter, op)
	}

	if !isValidUF(accessKey[0:2].String()) {
		return errors.E(ErrInvalidUF, ErrCodeInvalidUF, op)
	}

	if !isValidMonth(accessKey[4:6].String()) {
		return errors.E(ErrInvalidMonth, ErrCodeInvalidMonth, op)
	}

	if !isValidCPFCNPJ(accessKey[6:20].String()) {
		return errors.E(ErrInvalidCPFCNPJ, ErrCodeInvalidCPFCNPJ, op)
	}

	if !isValidModel(accessKey[20:22].String()) {
		return errors.E(ErrInvalidModel, ErrCodeInvalidModel, op)
	}

	if !isValidationDigitCorrect(accessKey.String()) {
		return errors.E(ErrInvalidDigit, ErrCodeInvalidDigit, op)
	}

	return nil
}

func isNFF(accessKey AccessKey) bool {
	/*model = 55 && tpEmis = 3 && AAMM more recent than April 2021*/
	if accessKey[20:22] == "55" && accessKey[34] == '3' && accessKey[2:6] >= "2104" {
		return true
	}
	return false
}

// validate execute all sub-routines necessary to perform a full accesskey validation for NFF
func validateNFF(accessKey AccessKey) error {
	const op errors.Op = "validateNFF"

	if !isValidSerieForNFF(accessKey[22:25].String()) {
		return errors.E(ErrInvalidSerieForNFF, ErrCodeInvalidSerieForNFF, op)
	}

	if !isValidNumeroForNFF(accessKey[25:34].String()) {
		return errors.E(ErrInvalidNumeroForNFF, ErrCodeInvalidNumeroForNFF, op)
	}

	if accessKey[29] == '1' {
		err := stakeholder.CheckCNPJ(accessKey[6:20].String())
		if err != nil {
			return errors.E(err, op, ErrCodeInvalidCNPJForNFF)
		}
	} else if accessKey[29] == '2' {
		if accessKey[6:9] != "000" {
			return errors.New("cpf is not padded with 0", ErrCodeInvalidCPFForNFF, op)
		}

		err := stakeholder.CheckCPF(accessKey[9:20].String())
		if err != nil {
			return errors.E(err, op, ErrCodeInvalidCPFForNFF)
		}
	}

	return nil
}

func isDigitOnly(accesskey AccessKey) bool {
	for _, token := range accesskey {
		if !(token >= '0' && token <= '9') {
			return false
		}
	}
	return true
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

func isValidMonth(month string) bool {
	if len(month) != 2 {
		return false
	}
	if month[0] == '0' {
		return month[1] != '0'
	}

	if month[0] == '1' {
		switch month[1] {
		case '0', '1', '2':
			return true
		default:
			return false
		}
	}
	return false
}

func isValidMonthDay(month string, day string) bool {
	if len(month) != 2 {
		return false
	}

	if len(day) != 2 {
		return false
	}

	if day == "00" {
		return false
	}

	switch month {
	case "01", "03", "05", "07", "08", "10", "12":
		if day > "31" {
			return false
		}
	case "04", "06", "09", "11":
		if day > "30" {
			return false
		}
	case "02":
		if day > "29" {
			return false
		}
	default:
		return false
	}

	return true
}

func isValidModel(model string) bool {
	switch model {
	case "1A",
		"01", "02",
		"04",
		"06", "07", "08", "09", "10", "11",
		"13", "14", "15", "16",
		"18",
		"21", "22",
		"26",
		"55",
		"57",
		"59", "60",
		"63",
		"65",
		"67":
		return true
	default:
		return false
	}
}

var validationDigitMultipliers = []int{
	4, 3, 2, 9, 8, 7, 6, 5, 4, 3,
	2, 9, 8, 7, 6, 5, 4, 3, 2, 9,
	8, 7, 6, 5, 4, 3, 2, 9, 8, 7,
	6, 5, 4, 3, 2, 9, 8, 7, 6, 5,
	4, 3, 2,
}

func isValidationDigitCorrect(accessKey string) bool {
	if len(accessKey) != 44 {
		return false
	}

	var sum int

	for i := 0; i < 43; i++ {
		sum += (int(accessKey[i]) - '0') * validationDigitMultipliers[i]
	}

	r := sum % 11
	vd := 0
	if r > 1 {
		vd = 11 - r
	}

	return (int(accessKey[43]) - '0') == vd
}

func isValidCPFCNPJ(cpfcnpj string) bool {
	if stakeholder.CheckCNPJ(cpfcnpj) == nil {
		return true
	}

	if cpfcnpj[0:3] != "000" {
		return false
	}

	return stakeholder.CheckCPF(cpfcnpj[3:14]) == nil
}

func isValidSerieForNFF(serie string) bool {
	return serie[0] != '0'
}

func isValidNumeroForNFF(numero string) bool {
	if !isValidMonthDay(numero[0:2], numero[2:4]) {
		return false
	}
	if numero[4] != '1' && numero[4] != '2' {
		return false
	}

	return true
}
