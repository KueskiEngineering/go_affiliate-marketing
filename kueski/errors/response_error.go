package errors

import "strconv"

// ResponseError - Type for custom API errors.
type ResponseError int

// InvalidCurp - Error for Invalid Curp
// InvalidEmail - Error for Invalid Email
// InvalidCurpAndEmail - Error for Invalid Curp And Email
// InvalidFullDataFormat - Error for Invalid Full Data Format
// InvalidRequestID - Error for Invalid Request Id
// InvalidFullDataAndRequestID - Error for Invalid Full Data And Request Id
// RequestIDNotFound - Error for Request Id Not Found
// MissingCurp - Error for Missing Curp
// MissingEmail - Error for Missing Email
// MissingCurpAndEmail - Error for Missing Curp And Email
// MissingCurpInvalidEmail - Error for Missing Curp Invalid Email
// MissingEmailInvalidCurp - Error for Missing Email Invalid Curp
// MissingFullData - Error for Missing Full Data
// MissingRequestID - Error for Missing Request Id
// MissingFullDataAndRequestID - Error for Missing Full Data And Request Id
// MissingFullDataInvalidRequestID - Error for Missing Full Data Invalid Request Id
// MissingRequestIDInvalidFullData - Error for Missing Request Id Invalid Full Data
// UnableToMakeConnection - Error for Unable To Make Connection
// LeadEvaluationMalformedRequest - Error for Lead Evaluation Malformed Request
// InvalidLeadEvaluationResponseFormat - Error for Invalid Lead Evaluation Response Format
// LeadDataMalformedRequest - Error for Lead Data Malformed Request
// InvalidLeadDataResponseFormat - Error for Invalid Lead Data Response Format
// AccessDenied - Error for Access Denied
// InvalidJWTResponseFormat - Error for Invalid Jwt Response Format
// InvalidExpirationFormat - Error for Invalid Expiration Format
// UnableToRefreshJWT - Error for Unable To Refresh Jwt
// InvalidSignatureFormat - Error for Invalid Signature Format
// ExpiredJWTToken - Error for Expired Jwt Token
// ExistingLead - Error for Existing Lead
// DuplicatedLead - Error for Duplicated Lead
// GeneralError - Error for General Error
const (
  InvalidCurp                 ResponseError = 1
  InvalidEmail                ResponseError = 2
  InvalidCurpAndEmail         ResponseError = 3
  InvalidFullDataFormat       ResponseError = 4
  InvalidRequestID            ResponseError = 5
  InvalidFullDataAndRequestID ResponseError = 6
  RequestIDNotFound           ResponseError = 7

  MissingCurp                     ResponseError = 11
  MissingEmail                    ResponseError = 12
  MissingCurpAndEmail             ResponseError = 13
  MissingCurpInvalidEmail         ResponseError = 14
  MissingEmailInvalidCurp         ResponseError = 15
  MissingFullData                 ResponseError = 16
  MissingRequestID                ResponseError = 17
  MissingFullDataAndRequestID     ResponseError = 18
  MissingFullDataInvalidRequestID ResponseError = 19
  MissingRequestIDInvalidFullData ResponseError = 20

  UnableToMakeConnection              ResponseError = 21
  LeadEvaluationMalformedRequest      ResponseError = 22
  InvalidLeadEvaluationResponseFormat ResponseError = 23
  LeadDataMalformedRequest            ResponseError = 24
  InvalidLeadDataResponseFormat       ResponseError = 25
  ErrorNotIdentifiedFromAPI           ResponseError = 26

  AccessDenied             ResponseError = 31
  InvalidJWTResponseFormat ResponseError = 32
  UnableToRefreshJWT       ResponseError = 33
  InvalidSignatureFormat   ResponseError = 34
  ExpiredJWTToken          ResponseError = 35

  ExistingLead   ResponseError = 41
  DuplicatedLead ResponseError = 42

  GeneralError ResponseError = 99
)

type errorDescription struct {
  name string
  desc string
}

var descriptions = map[ResponseError]errorDescription{
  InvalidCurp:                         errorDescription{"InvalidCurp", "Invalid CURP."},
  InvalidEmail:                        errorDescription{"InvalidEmail", "Invalid e-mail."},
  InvalidCurpAndEmail:                 errorDescription{"InvalidCurpAndEmail", "Invalid CURP And e-mail."},
  InvalidFullDataFormat:               errorDescription{"InvalidFullDataFormat", "Invalid full data format."},
  InvalidRequestID:                    errorDescription{"InvalidRequestID", "Invalid request ID."},
  InvalidFullDataAndRequestID:         errorDescription{"InvalidFullDataAndRequestID", "Invalid full data And request ID."},
  RequestIDNotFound:                   errorDescription{"RequestIDNotFound", "Request ID Not Found."},
  MissingCurp:                         errorDescription{"MissingCurp", "Missing CURP."},
  MissingEmail:                        errorDescription{"MissingEmail", "Missing e-mail."},
  MissingCurpAndEmail:                 errorDescription{"MissingCurpAndEmail", "Missing CURP And e-mail."},
  MissingCurpInvalidEmail:             errorDescription{"MissingCurpInvalidEmail", "Missing CURP and Invalid e-mail."},
  MissingEmailInvalidCurp:             errorDescription{"MissingEmailInvalidCurp", "Missing e-mail and Invalid CURP."},
  MissingFullData:                     errorDescription{"MissingFullData", "Missing full data."},
  MissingRequestID:                    errorDescription{"MissingRequestID", "Missing request ID."},
  MissingFullDataAndRequestID:         errorDescription{"MissingFullDataAndRequestID", "Missing full data and request ID."},
  MissingFullDataInvalidRequestID:     errorDescription{"MissingFullDataInvalidRequestID", "Missing full data and Invalid request ID."},
  MissingRequestIDInvalidFullData:     errorDescription{"MissingRequestIDInvalidFullData", "Missing request ID and Invalid full data."},
  UnableToMakeConnection:              errorDescription{"UnableToMakeConnection", "Unable to make connection."},
  LeadEvaluationMalformedRequest:      errorDescription{"LeadEvaluationMalformedRequest", "Lead Evaluation malformed request."},
  InvalidLeadEvaluationResponseFormat: errorDescription{"InvalidLeadEvaluationResponseFormat", "Invalid Lead Evaluation response format."},
  ErrorNotIdentifiedFromAPI:           errorDescription{"ErrorNotIdentifiedFromAPI", "Error from the API is not recognized"},
  LeadDataMalformedRequest:            errorDescription{"LeadDataMalformedRequest", "Lead Data malformed request."},
  InvalidLeadDataResponseFormat:       errorDescription{"InvalidLeadDataResponseFormat", "Invalid Lead Data response format."},
  AccessDenied:                        errorDescription{"AccessDenied", "Access denied."},
  InvalidJWTResponseFormat:            errorDescription{"InvalidJWTResponseFormat", "Invalid JWT response format."},
  UnableToRefreshJWT:                  errorDescription{"UnableToRefreshJWT", "Unable to refresh JWT."},
  InvalidSignatureFormat:              errorDescription{"InvalidSignatureFormat", "Invalid Signature format."},
  ExpiredJWTToken:                     errorDescription{"ExpiredJWTToken", "Expired JWT token."},
  ExistingLead:                        errorDescription{"ExistingLead", "Existing Lead."},
  DuplicatedLead:                      errorDescription{"DuplicatedLead", "Duplicated Lead."},
  GeneralError:                        errorDescription{"GeneralError", "General error."},
}

func (i ResponseError) String() string {
  return i.retrieveDescription().name
}

func (i ResponseError) Error() string {
  return i.String()
}

// Description - Response Error long description.
func (i ResponseError) Description() string {
  return i.retrieveDescription().desc
}

// StringToResponseError - Convert a string into a ResponseError. If it doesn't match, returns nil.
func StringToResponseError(name string) error {
  var responseError error
  responseError = nil

  for errorKey, description := range descriptions {
    if description.name == name {
      responseError = errorKey
      break
    }
  }

  return responseError
}

func (i ResponseError) retrieveDescription() errorDescription {
  description := descriptions[i]
  if description == (errorDescription{}) {
    description.name = "ResponseError(" + strconv.FormatInt(int64(i), 10) + ")"
    description.desc = "Unknown response error."
  }
  return description
}
