package kueski

import (
  "encoding/json"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/util"
)

type evaluation struct {
  Curp  string `json:"curp"`
  Email string `json:"email"`
}

type evaluationResponse struct {
  Email     string
  Curp      string
  Status    string
  RequestID string `json:"request_id"`
}

// leadEvaluationPath - Path to the lead evaluation with CURP/email endpoint.
const leadEvaluationPath string = "affiliates/lead-evaluation"

var evaluationErrors = map[string]error{
  "email is invalid, curp is missing": errors.MissingCurpInvalidEmail,
  "email is invalid, curp is invalid": errors.InvalidCurpAndEmail,
  "email is invalid":                  errors.InvalidEmail,
  "email is missing, curp is missing": errors.MissingCurpAndEmail,
  "email is missing, curp is invalid": errors.MissingEmailInvalidCurp,
  "email is missing":                  errors.MissingEmail,
  "curp is missing":                   errors.MissingCurp,
  "curp is invalid":                   errors.InvalidCurp,
}

var validResponseStatus = map[string]bool{"approved": true, "duplicated": true, "existing": true}

// leadEvaluator - Type that defines the lead evaluation signature.
type leadEvaluator func(client *Client, curp, email string) (string, error)

func leadEvaluation(client *Client, curp, email string) (string, error) {
  body, _ := json.Marshal(evaluation{curp, email})
  response, err := client.makeRequest(leadEvaluationPath, body)

  if err != nil {
    return "", err
  }

  responseBody, err := util.ExtractBody(response)

  if err != nil {
    return "", errors.InvalidLeadEvaluationResponseFormat
  }

  if response.StatusCode == 500 {
    return "", errors.LeadEvaluationMalformedRequest
  }

  if response.StatusCode == 400 {
    return "", resolveAPIError(responseBody, errors.InvalidLeadEvaluationResponseFormat, evaluationErrors)
  }

  if response.StatusCode == 401 {
    return "", errors.ExpiredJWTToken
  }

  return resolveEvaluationResponse(responseBody, curp, email)
}

func resolveEvaluationResponse(body []byte, curp, email string) (string, error) {
  var response evaluationResponse
  responseErr := json.Unmarshal(body, &response)

  if responseErr != nil || !(response.RequestID != "" && response.Curp == curp && response.Email == email && validResponseStatus[response.Status]) {
    return "", errors.InvalidLeadEvaluationResponseFormat
  }

  var err error

  if response.Status == "duplicated" {
    err = errors.DuplicatedLead
  }

  if response.Status == "existing" {
    err = errors.ExistingLead
  }

  return response.RequestID, err
}
