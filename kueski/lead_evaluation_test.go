package kueski

import (
  "fmt"
  "net/http"
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

var curp = "CURP2345678"
var email = "e@mail"
var requestID = "456789"

func TestLeadEvaluation(t *testing.T) {
  client := Client{}
  client.jwtProvider = &fakeTokenProvider{true}

  responses := []*http.Response{
    invalidHTTPResponse(),
    buildHTTPResponse(500, ""),
    buildHTTPResponse(401, ""),
    buildHTTPResponse(400, `{ "error": "email is invalid, curp is missing" }`),
    buildHTTPResponse(400, `{ "error": "email is invalid, curp is invalid" }`),
    buildHTTPResponse(400, `{ "error": "email is invalid" }`),
    buildHTTPResponse(400, `{ "error": "email is missing, curp is missing" }`),
    buildHTTPResponse(400, `{ "error": "email is missing, curp is invalid" }`),
    buildHTTPResponse(400, `{ "error": "email is missing" }`),
    buildHTTPResponse(400, `{ "error": "curp is missing" }`),
    buildHTTPResponse(400, `{ "error": "curp is invalid" }`),
    buildHTTPResponse(201, fmt.Sprintf(`{ "curp": "%s", "email": "%s", "request_id": "%s", "status": "approved" }`, curp, email, requestID)),
  }
  errors := []error{
    errors.InvalidLeadEvaluationResponseFormat,
    errors.LeadEvaluationMalformedRequest,
    errors.ExpiredJWTToken,
    errors.MissingCurpInvalidEmail,
    errors.InvalidCurpAndEmail,
    errors.InvalidEmail,
    errors.MissingCurpAndEmail,
    errors.MissingEmailInvalidCurp,
    errors.MissingEmail,
    errors.MissingCurp,
    errors.InvalidCurp,
    nil,
  }

  for i, response := range responses {
    err := errors[i]
    requester := func(url string, headers map[string]string, body []byte) (*http.Response, error) { return response, nil }

    client.requester = requester
    responseGot, dataErr := leadEvaluation(&client, curp, email)

    expected := ""
    if i == len(responses)-1 {
      expected = requestID
    }

    assert.Equal(t, expected, responseGot, fmt.Sprintf("Test case %d", i))
    assert.Equal(t, err, dataErr)
  }
}

func TestLeadEvaluationRequestError(t *testing.T) {
  client := Client{}
  client.jwtProvider = &fakeTokenProvider{false}

  _, err := leadEvaluation(&client, curp, email)
  assert.Equal(t, errors.GeneralError, err)
}

func TestResolveEvaluationResponse(t *testing.T) {
  bodys := []string{
    `{ "curp": "CURP123", "email:`,
    fmt.Sprintf(`{ "curp": "%s", "email": "%s" }`, curp, email),
    fmt.Sprintf(`{ "curp": "%s", "email": "%s", "request_id": "%s", "status": "other" }`, curp, email, requestID),
    fmt.Sprintf(`{ "curp": "OtherCURP", "email": "%s", "request_id": "%s", "status": "approved" }`, email, requestID),
    fmt.Sprintf(`{ "curp": "%s", "email": "OtherEmail", "request_id": "%s", "status": "approved" }`, curp, requestID),
    fmt.Sprintf(`{ "curp": "%s", "email": "%s", "request_id": "%s", "status": "duplicated" }`, curp, email, requestID),
    fmt.Sprintf(`{ "curp": "%s", "email": "%s", "request_id": "%s", "status": "existing" }`, curp, email, requestID),
    fmt.Sprintf(`{ "curp": "%s", "email": "%s", "request_id": "%s", "status": "approved" }`, curp, email, requestID),
  }
  responses := []string{"", "", "", "", "", requestID, requestID, requestID}
  errors := []error{
    errors.InvalidLeadEvaluationResponseFormat,
    errors.InvalidLeadEvaluationResponseFormat,
    errors.InvalidLeadEvaluationResponseFormat,
    errors.InvalidLeadEvaluationResponseFormat,
    errors.InvalidLeadEvaluationResponseFormat,
    errors.DuplicatedLead,
    errors.ExistingLead,
    nil,
  }

  for i, body := range bodys {
    returnedRequestID, err := resolveEvaluationResponse([]byte(body), curp, email)
    assert.Equal(t, responses[i], returnedRequestID, fmt.Sprintf("Test case %d", i))
    assert.Equal(t, errors[i], err)
  }
}
