package accesskey

import (
	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/sefaz/stakeholder"
)

// validate execute all sub-routines necessary to perform a full accesskey validation
func validate(accessKey AccessKey) error {
	const op errors.Op = "validate"
	if accessKey == "" {
		return errors.E(op, ErrEmptyAccessKey, ErrCodeEmptyAccessKey)
	}
	if len(accessKey) != 44 {
		return errors.E(op, ErrInvalidLenght, ErrCodeInvalidAccessKey)
	}

	if !isDigitOnly(accessKey) {
		return errors.E(op, ErrInvalidCharacter, ErrCodeInvalidAccessKey)
	}

	if !isValidUF(accessKey[0:2].String()) {
		return errors.E(op, ErrInvalidUF, ErrCodeInvalidAccessKey)
	}

	if !isValidMonth(accessKey[4:6].String()) {
		return errors.E(op, ErrInvalidMonth, ErrCodeInvalidAccessKey)
	}

	if !isValidCPFCNPJ(accessKey[6:20].String()) {
		return errors.E(op, ErrInvalidCPFCNPJ, ErrCodeInvalidAccessKey)
	}

	if !isValidModel(accessKey[20:22].String()) {
		return errors.E(op, ErrInvalidModel, ErrCodeInvalidAccessKey)
	}

	if !isValidationDigitCorrect(accessKey.String()) {
		return errors.E(op, ErrInvalidDigit, ErrCodeInvalidAccessKey)
	}

	return nil
}

func (v *validator) Check(accessKey AccessKey) error {
	return validate(accessKey)
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
	} else {
		return false
	}
}

func isValidModel(model string) bool {
	return model == "01" || model == "1A" ||
		model == "02" || model == "04" ||
		model == "06" || model == "07" ||
		model == "08" || model == "09" ||
		model == "10" || model == "11" ||
		model == "13" || model == "14" ||
		model == "15" || model == "16" ||
		model == "18" || model == "21" ||
		model == "22" || model == "26" ||
		model == "55" || model == "57" ||
		model == "59" || model == "60" ||
		model == "63" || model == "65" ||
		model == "67"
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
