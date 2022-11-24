package lib_utils

// Error responses
// errorXXX occurs when XXX
const (
	ErrorParseJws                 = `error parsing into JWS`
	ErrorParseX509                = `error parsing into X509`
	ErrorBase64                   = `error decoding into base64`
	ErrorVerifying                = `error verifying signature`
	ErrorInvalidSpecificOperation = "invalid operation in state %s: %s"
	ErrorInvalidOperation         = "invalid operation in state"
	ErrorFailedWorldState         = "failed to get in state %s: %s"
	ErrorWorldState               = "unable to interact with world state. %s"
	ErrorNotExistInState          = "no state found for %s"
	ErrorAlreadyExistInState      = "the value %s already exists in the state database"
	ErrorIDSame                   = "the ids are the same"
	ErrorUnmarshal                = "unmarshal error %s"
	ErrorMarshal                  = "marshal error %s"
	ErrorGenerateKey              = "error generating key"
)

// Each code must be 4 characters

const (
	CodCert = "CERT"
)

// contract name
const (
	ContractNameCommon      = "common"
	ContractNameCertificate = "certificate"
)
