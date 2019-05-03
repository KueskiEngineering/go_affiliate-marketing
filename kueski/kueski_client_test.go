package kueski

import (
  "net/http"
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestResolveAPIError(t *testing.T) {
  malformeds := [2]string{"{", "{ \"error\": \"error\""}
  unknowns := [2]string{"{ \"errorf\": \"error\" }", "{ \"herror\": \"error\" }"}

  malformedError := errors.GeneralError
  errorMap := map[string]error{"error": errors.InvalidCurp}

  for _, str := range malformeds {
    err := resolveAPIError([]byte(str), malformedError, errorMap)
    assert.Equal(t, malformedError, err)
  }

  for _, str := range unknowns {
    err := resolveAPIError([]byte(str), malformedError, errorMap)
    assert.Equal(t, errors.ErrorNotIdentifiedFromAPI, err)
  }

  assert.Equal(t, errors.InvalidCurp, resolveAPIError([]byte("{ \"error\": \"error\" }"), malformedError, errorMap))
}

func TestMakeRequestErrors(t *testing.T) {
  client := Client{}
  client.jwtProvider = &fakeTokenProvider{false}

  response, err := client.makeRequest("path", []byte("body"))

  assert.Nil(t, response)
  assert.Equal(t, errors.GeneralError, err)
}

func TestMakeRequest(t *testing.T) {
  requestBody := []byte("Body")
  requester := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    assert.Equal(t, "http://kueski.com/path", url)
    assert.Equal(t, 2, len(headers))
    assert.Equal(t, headers[Authorization], "Bearer Token")
    assert.Equal(t, headers[ContentType], ApplicationJSON)
    assert.Equal(t, body, requestBody)

    return nil, errors.GeneralError
  }

  client := Client{}
  client.jwtProvider = &fakeTokenProvider{true}
  client.url = "http://kueski.com"
  client.requester = requester

  response, err := client.makeRequest("path", []byte("Body"))

  assert.Nil(t, response)
  assert.Equal(t, errors.GeneralError, err)
}

func TestEvaluate(t *testing.T) {
  validationWorks := false
  evaluationWorks := false
  dataWorks := false

  testCurp := "CURP"
  testEmail := "e@mail"
  testData := fakeTokenProvider{false}
  testRequestID := "abc34567"

  testClient := Client{}

  fakeValidator := func(curp, email string, fullData interface{}) error {
    assert.Equal(t, testCurp, curp)
    assert.Equal(t, testEmail, email)
    assert.Equal(t, testData, fullData)

    if validationWorks {
      return nil
    }

    return errors.InvalidCurpAndEmail
  }

  fakeEvaluator := func(client *Client, curp, email string) (string, error) {
    assert.Equal(t, &testClient, client)
    assert.Equal(t, testCurp, curp)
    assert.Equal(t, testEmail, email)

    if evaluationWorks {
      return testRequestID, nil
    }

    return testRequestID, errors.DuplicatedLead
  }

  fakeDataHandler := func(client *Client, jsonData interface{}, requestID string) error {
    assert.Equal(t, &testClient, client)
    assert.Equal(t, testData, jsonData)
    assert.Equal(t, testRequestID, requestID)

    if dataWorks {
      return nil
    }

    return errors.RequestIDNotFound
  }

  testClient.validator = fakeValidator
  testClient.evaluator = fakeEvaluator
  testClient.dataHandler = fakeDataHandler

  returnedRequestID, err := testClient.Evaluate(testCurp, testEmail, testData)

  assert.Equal(t, "", returnedRequestID)
  assert.Equal(t, errors.InvalidCurpAndEmail, err)

  validationWorks = true
  returnedRequestID, err = testClient.Evaluate(testCurp, testEmail, testData)

  assert.Equal(t, testRequestID, returnedRequestID)
  assert.Equal(t, errors.DuplicatedLead, err)

  evaluationWorks = true
  returnedRequestID, err = testClient.Evaluate(testCurp, testEmail, testData)

  assert.Equal(t, testRequestID, returnedRequestID)
  assert.Equal(t, errors.RequestIDNotFound, err)

  dataWorks = true
  returnedRequestID, err = testClient.Evaluate(testCurp, testEmail, testData)

  assert.Equal(t, testRequestID, returnedRequestID)
  assert.Nil(t, err)
}
