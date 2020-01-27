package accesskey

import (
	"strings"

	"github.com/arquivei/foundationkit/errors"
)

// validate execute all sub-routines necessary to perform a full accesskey validation
func validate(accessKey Key) error {
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

	if !isValidUF(accessKey[0:2]) {
		return errors.E(op, ErrInvalidUF, ErrCodeInvalidAccessKey)
	}

	if !isValidMonth(accessKey[4:6]) {
		return errors.E(op, ErrInvalidMonth, ErrCodeInvalidAccessKey)
	}

	if !isValidCPFCNPJ(accessKey[6:20].String()) {
		return errors.E(op, ErrInvalidCPFCNPJ, ErrCodeInvalidAccessKey)
	}

	if !isValidModel(accessKey[20:22]) {
		return errors.E(op, ErrInvalidModel, ErrCodeInvalidAccessKey)
	}

	if !isValidationDigitCorrect(accessKey) {
		return errors.E(op, ErrInvalidDigit, ErrCodeInvalidAccessKey)
	}

	return nil
}

func (v *validator) Check(accessKey Key) error {
	return validate(accessKey)
}

func isDigitOnly(accesskey Key) bool {
	for _, token := range accesskey {
		if !(token >= '0' && token <= '9') {
			return false
		}
	}
	return true
}

func isValidUF(uf Key) bool {
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

func isValidMonth(month Key) bool {
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

func isValidModel(model Key) bool {
	return model == "55"
}

var validationDigitMultipliers = []int{
	4, 3, 2, 9, 8, 7, 6, 5, 4, 3,
	2, 9, 8, 7, 6, 5, 4, 3, 2, 9,
	8, 7, 6, 5, 4, 3, 2, 9, 8, 7,
	6, 5, 4, 3, 2, 9, 8, 7, 6, 5,
	4, 3, 2,
}

func isValidationDigitCorrect(accessKey Key) bool {
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
	return isValidCNPJ(cpfcnpj) || isValidCPF(cpfcnpj)
}

func isValidCPF(cpfcnpj string) bool {
	if len(cpfcnpj) != 14 {
		return false
	}
	if !strings.HasPrefix(cpfcnpj, "000") {
		return false
	}
	cpf := cpfcnpj[3:14]

	switch cpf {
	case
		"00000000000",
		"11111111111",
		"22222222222",
		"33333333333",
		"44444444444",
		"55555555555",
		"66666666666",
		"77777777777",
		"88888888888",
		"99999999999":
		return false
	}

	type validateData struct {
		multiplier1 uint16
		multiplier2 uint16
	}

	data := [11]validateData{
		{10, 11}, {9, 10}, {8, 9}, {7, 8}, {6, 7}, {5, 6},
		{4, 5}, {3, 4}, {2, 3}, {0, 2}, {0, 0},
	}

	sumDigit1 := uint16(0)
	sumDigit2 := uint16(0)
	for i := 0; i < 11; i++ {
		c := uint16(cpf[i]) - '0'
		if c > 9 {
			return false
		}

		sumDigit1 += c * data[i].multiplier1
		sumDigit2 += c * data[i].multiplier2
	}

	if uint16(cpf[9]-'0') != (sumDigit1*10%11%10) || uint16(cpf[10]-'0') != (sumDigit2*10%11%10) {
		return false
	}

	return true
}

func isValidCNPJ(cnpj string) bool {
	cnpjLength := len(cnpj)
	if cnpjLength == 0 {
		return false
	}
	if len(cnpj) != 14 {
		return false
	}

	type validateData struct {
		multiplier1 uint16
		multiplier2 uint16
	}

	data := [14]validateData{
		{6, 5}, {7, 6}, {8, 7}, {9, 8}, {2, 9}, {3, 2},
		{4, 3}, {5, 4}, {6, 5}, {7, 6}, {8, 7}, {9, 8},
		{0, 9}, {0, 0},
	}

	sumDigit1 := uint16(0)
	sumDigit2 := uint16(0)
	for i := 0; i < 14; i++ {
		c := uint16(cnpj[i]) - '0'
		if c > 9 {
			return false
		}

		sumDigit1 += c * data[i].multiplier1
		sumDigit2 += c * data[i].multiplier2
	}

	if uint16(cnpj[12]-'0') != (sumDigit1%11%10) || uint16(cnpj[13]-'0') != (sumDigit2%11%10) {
		return false
	}

	return true
}

/* isValidEmissionType
This value comes from the field: tpEmis
According to "Manual de Orientação do Contribuinte":

Tipo de Emissão da NF-e
1 - Emissão normal (não em contingência);
2 - Contingência FS-IA, com impressão do DANFE em formulário de segurança;
3 - Contingência SCAN (Sistema de Contingência do Ambiente Nacional);
4 - Contingência DPEC (Declaração Prévia da Emissão em Contingência);
5 - Contingência FS-DA, com impressão do DANFE em formulário de segurança;
6 - Contingência SVC-AN (SEFAZ Virtual de Contingência do AN);
7 - Contingência SVC-RS (SEFAZ Virtual de Contingência do RS);
9 - Contingência off-line da NFC-e (as demais opções de contingência são válidas também para a NFC-e). Para a NFC-e somente estão disponíveis e são válidas as opções de contingência 5 e 9.

ps: this validation is only valid for accesskeys of
V2 or more, as this field does not appear on V1 access keys.
*/
