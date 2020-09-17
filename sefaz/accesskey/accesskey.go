package accesskey

// AccessKey represents the AccessKey entity
type AccessKey string

// Validator is the interface that performs AccessKey validation
// Deprecated: prefer using the Check function
type Validator interface {
	// Deprecated: prefer using the Check function
	Check(AccessKey) error
}

type validator struct{}

// NewValidator returns a real validator
// Deprecated: prefer using the Check function instead of instantiating a Validator
func NewValidator() Validator {
	return &validator{}
}

func (a AccessKey) String() string {
	return string(a)
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
9 - Contingência off-line da NFC-e (as demais opções de contingência são válidas também para a NFC-e).
Para a NFC-e somente estão disponíveis e são válidas as opções de contingência 5 e 9.

PS: This validation must be for access keys V2 or more, as this field does not appear on V1 access keys.
*/
