package stakeholder

import "strconv"

func isDigitOnly(accesskey string) bool {
	for _, token := range accesskey {
		if token < '0' || token > '9' {
			return false
		}
	}
	return true
}

// CheckCPF validates the given CPF
func CheckCPF(cpf string) error {
	if cpf == "" {
		return ErrCPFEmpty
	}

	if !isDigitOnly(cpf) {
		return ErrCPFNotDigitsOnly
	}

	if len(cpf) != 11 {
		return ErrCPFInvalidLength
	}

	i, sum := 10, 0

	for index := 0; index < len(cpf)-2; index++ {
		pos, err := strconv.Atoi(string(cpf[index]))
		if err != nil {
			return err
		}

		sum += pos * i
		i--
	}

	prod := sum * 10
	mod := prod % 11

	if mod == 10 {
		mod = 0
	}

	digit1, err := strconv.Atoi(string(cpf[9]))
	if err != nil {
		return err
	}

	if mod != digit1 {
		return ErrCPFValidationDigit
	}

	i, sum = 11, 0

	for index := 0; index < len(cpf)-1; index++ {
		pos, err := strconv.Atoi(string(cpf[index]))
		if err != nil {
			return err
		}

		sum += pos * i
		i--
	}

	prod = sum * 10
	mod = prod % 11

	if mod == 10 {
		mod = 0
	}

	digit2, err := strconv.Atoi(string(cpf[10]))
	if err != nil {
		return err
	}

	if mod != digit2 {
		return ErrCPFValidationDigit
	}

	return nil
}

// CheckCNPJ validates the given CNPJ
func CheckCNPJ(cnpj string) error {
	if len(cnpj) == 0 {
		return ErrCNPJEmpty
	}

	if !isDigitOnly(cnpj) {
		return ErrCNPJNotDigitsOnly
	}

	if len(cnpj) != 14 {
		return ErrCNPJInvalidLength
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
			return ErrCNPJNotDigitsOnly
		}

		sumDigit1 += c * data[i].multiplier1
		sumDigit2 += c * data[i].multiplier2
	}

	if uint16(cnpj[12]-'0') != (sumDigit1%11%10) || uint16(cnpj[13]-'0') != (sumDigit2%11%10) {
		return ErrCNPJValidationDigit
	}

	return nil
}
