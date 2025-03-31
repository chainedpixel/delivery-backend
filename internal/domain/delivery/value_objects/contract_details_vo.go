package value_objects

import (
	"encoding/json"
	"time"
)

// ContractDetails representa los detalles del contrato de una empresa
type ContractDetails struct {
	ContractType   string    `json:"contract_type"`
	PaymentTerms   string    `json:"payment_terms"`
	RenewalType    string    `json:"renewal_type"`
	NoticePeriod   int       `json:"notice_period"` // días
	SpecialClauses []string  `json:"special_clauses,omitempty"`
	SignedBy       string    `json:"signed_by,omitempty"`
	SignedAt       time.Time `json:"signed_at,omitempty"`
}

func NewContractDetails(contractType, paymentTerms, renewalType string, noticePeriod int) *ContractDetails {
	return &ContractDetails{
		ContractType: contractType,
		PaymentTerms: paymentTerms,
		RenewalType:  renewalType,
		NoticePeriod: noticePeriod,
	}
}

func (cd *ContractDetails) WithSpecialClauses(clauses []string) *ContractDetails {
	cd.SpecialClauses = clauses
	return cd
}

func (cd *ContractDetails) WithSignature(signedBy string, signedAt time.Time) *ContractDetails {
	cd.SignedBy = signedBy
	cd.SignedAt = signedAt
	return cd
}

// IsValid valida los detalles del contrato
func (cd *ContractDetails) IsValid() bool {
	return cd.ContractType != "" &&
		cd.PaymentTerms != "" &&
		cd.RenewalType != "" &&
		cd.NoticePeriod >= 0
}

// ToJSON convierte los detalles del contrato a un string JSON
func (cd *ContractDetails) ToJSON() (string, error) {
	data, err := json.Marshal(cd)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON crea un objeto ContractDetails desde un string JSON
func ContractDetailsFromJSON(jsonStr string) (*ContractDetails, error) {
	var details ContractDetails
	if err := json.Unmarshal([]byte(jsonStr), &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// ToString convierte los detalles del contrato a un string legible
func (cd *ContractDetails) ToString() string {
	json, _ := cd.ToJSON()
	return json
}

// Equals compara si dos objetos ContractDetails son iguales
func (cd *ContractDetails) Equals(value ValidaterObject[ContractDetails]) bool {
	other := value.GetValue()

	// Comparación básica de atributos principales
	if cd.ContractType != other.ContractType ||
		cd.PaymentTerms != other.PaymentTerms ||
		cd.RenewalType != other.RenewalType ||
		cd.NoticePeriod != other.NoticePeriod ||
		cd.SignedBy != other.SignedBy {
		return false
	}

	// Comparación de SignedAt (si existe)
	if (!cd.SignedAt.IsZero() || !other.SignedAt.IsZero()) &&
		!cd.SignedAt.Equal(other.SignedAt) {
		return false
	}

	// Comparación de SpecialClauses
	if len(cd.SpecialClauses) != len(other.SpecialClauses) {
		return false
	}

	for i, clause := range cd.SpecialClauses {
		if i >= len(other.SpecialClauses) || clause != other.SpecialClauses[i] {
			return false
		}
	}

	return true
}

// GetValue retorna el valor del objeto
func (cd *ContractDetails) GetValue() ContractDetails {
	return *cd
}
