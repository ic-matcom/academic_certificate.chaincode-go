package certificate

type StateValidation uint

const (
	Invalid  StateValidation = iota // invalidated for some reason
	New                             // certificate without signatures
	SignedS                         // signed by Secretary
	SignedSD                        // signed by Secretary and Dean
	Valid                           // signed by Secretary, Dean and Rector
)

type ValidatorType uint

const (
	NoValidator ValidatorType = iota // invalidated for some reason
	Secretary                        // certificate without signatures
	Dean                             // signed by Secretary
	Rector                           // signed by Secretary and Dean
)

//Auxiliary Functions
func (state StateValidation) String() string {
	names := []string{"Invalid", "Miss Dean and Rector signatures", "Miss Rector signature", "Va"}
	if state < Invalid || state > Valid {
		return "unknown"
	}
	return names[state]
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType               string          `json:"docType"`
	ID                    string          `json:"ID"`
	Certification         string          `json:"certification"`
	GoldCertificate       bool            `json:"gold_certificate"`
	Emitter               string          `json:"emitter"`
	Accredited            string          `json:"accredited"`
	Date                  string          `json:"date"`
	CreatedBy             string          `json:"created_by"`
	SecretaryValidating   string          `json:"secretary_validating"`
	DeanValidating        string          `json:"dean_validating"`
	RectorValidating      string          `json:"rector_validating"`
	FacultyVolumeFolio    string          `json:"volume_folio_faculty"`
	UniversityVolumeFolio string          `json:"volume_folio_university"`
	InvalidReason         string          `json:"invalid_reason"`
	Status                StateValidation `json:"certificate_status"`
}

type GetRequest struct {
	ID string `json:"id"`
}

type ValidateAsset struct {
	ID         string        `json:"ID"`
	Validator  string        `json:"validator"`
	ValidatorT ValidatorType `json:"validator_type"`
}

type InvalidateAsset struct {
	ID          string `json:"ID"`
	Description string `json:"description"`
}
