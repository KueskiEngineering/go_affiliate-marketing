package kueski

import (
  "fmt"
  "net/http"
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestResolveLeadDataResponse(t *testing.T) {
  requestID := "12345678"
  good := fmt.Sprintf(`{ "response": "ok", "request_id": "%s" }`, requestID)
  bads := []string{
    fmt.Sprintf(`{ "response": "ok", "request_id": "%s`, requestID),
    fmt.Sprintf(`{ "other": "ok", "keys": "%s" }`, requestID),
    `{ "response": "ok", "request_id": "91919191" }`,
    fmt.Sprintf(`{ "response": "other", "request_id": "%s" }`, requestID),
  }

  assert.Nil(t, resolveLeadDataResponse([]byte(good), requestID))

  for _, bad := range bads {
    assert.Equal(t, errors.InvalidLeadDataResponseFormat, resolveLeadDataResponse([]byte(bad), requestID))
  }
}

type sampleData struct {
  name string
  ssid int
}

func TestLeadDataHandlerRequestError(t *testing.T) {
  client := Client{}
  client.jwtProvider = &fakeTokenProvider{false}
  data := sampleData{"Fake Name", 1}
  requestID := "0987654321"

  err := leadData(&client, data, requestID)
  assert.Equal(t, errors.GeneralError, err)
}

func TestLeadDataResponseErrors(t *testing.T) {
  client := Client{}
  client.jwtProvider = &fakeTokenProvider{true}

  data := sampleData{"Fake Name", 1}
  requestID := "0987654321"

  responses := []*http.Response{
    invalidHTTPResponse(),
    buildHTTPResponse(500, ""),
    buildHTTPResponse(400, `{ "error": "Request not found" }`),
    buildHTTPResponse(400, `{ "error": "request_id is invalid, full_data is missing" }`),
    buildHTTPResponse(400, `{ "error": "request_id is invalid, full_data is invalid" }`),
    buildHTTPResponse(400, `{ "error": "request_id is invalid" }`),
    buildHTTPResponse(400, `{ "error": "request_id is missing, full_data is missing" }`),
    buildHTTPResponse(400, `{ "error": "request_id is missing, full_data is invalid" }`),
    buildHTTPResponse(400, `{ "error": "request_id is missing" }`),
    buildHTTPResponse(400, `{ "error": "full_data is missing" }`),
    buildHTTPResponse(400, `{ "error": "full_data is invalid" }`),
    buildHTTPResponse(401, ``),
    buildHTTPResponse(201, fmt.Sprintf(`{ "response": "ok", "request_id": "%s" }`, requestID)),
  }

  errors := []error{
    errors.InvalidLeadDataResponseFormat,
    errors.LeadDataMalformedRequest,
    errors.RequestIDNotFound,
    errors.MissingFullDataInvalidRequestID,
    errors.InvalidFullDataAndRequestID,
    errors.InvalidRequestID,
    errors.MissingFullDataAndRequestID,
    errors.MissingRequestIDInvalidFullData,
    errors.MissingRequestID,
    errors.MissingFullData,
    errors.InvalidFullDataFormat,
    errors.ExpiredJWTToken,
    nil,
  }

  for i, response := range responses {
    err := errors[i]
    requester := func(url string, headers map[string]string, body []byte) (*http.Response, error) { return response, nil }

    client.requester = requester
    dataErr := leadData(&client, data, requestID)

    assert.Equal(t, err, dataErr)
  }
}
