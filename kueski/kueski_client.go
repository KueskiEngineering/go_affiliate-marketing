package kueski

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/util"
)

type leadValidator func(curp, email string, fullData interface{}) error

// Client - Interface to connect with Kueski Affiliates API.
type Client struct {
  url         string
  apiKey      string
  secretKey   string
  requester   util.PostRequestFunc
  validator   leadValidator
  evaluator   leadEvaluator
  dataHandler leadDataHandler
  jwtProvider TokenProvider
}

type apiError struct {
  Error string
}

const (
  tokenHeaderFormat string = "Bearer %s"
)

// NewClient - Constructor for Kueski API Client.
// url - API host URL.
// apiKey - API key.
// secretKey - Service secret key.
func NewClient(url, apiKey, secretKey string) *Client {
  client := new(Client)
  client.url = url
  client.apiKey = apiKey
  client.secretKey = secretKey
  client.requester = util.PostRequest
  client.evaluator = leadEvaluation
  client.dataHandler = leadData
  client.jwtProvider = NewJWTProvider()

  validator := Validator{util.ValidateCurp, util.ValidateEmail, util.ValidateFullData}
  client.validator = validator.validate

  return client
}

// Evaluate - Performs the lead evaluation.
// curp - Lead CURP.
// email - Lead email.
// fullData - Struct with the full lead data to be sent to Kueski API.
// Returns: A string with the Request ID if the operation was successful, otherwise error will not be nil.
// --
// Steps to follow:
// * Request a JWT.
// * Call LeadEvaluation with CURP and email.
// * Identify if response is successful, otherwise return proper error code.
// * Request a JWT.
// * Call LeadData with fullData and Request ID.
// * Identify if response is successful, otherwise return proper error code.
// * Return Request ID if all successful.
func (client *Client) Evaluate(curp, email string, fullData interface{}) (string, error) {
  // Validate data to POST before calling the API.
  err := client.validator(curp, email, fullData)

  if err != nil {
    return "", err
  }

  // Calls to the Kueski API.
  requestID, err := client.evaluator(client, curp, email)

  if err != nil {
    return requestID, err
  }

  err = client.dataHandler(client, fullData, requestID)

  if err != nil {
    return requestID, err
  }

  return requestID, nil
}

func (client *Client) makeRequest(path string, body []byte) (*http.Response, error) {
  token, err := client.jwtProvider.Token(client)

  if err != nil {
    return nil, err
  }

  url := util.BuildURL(client.url, path)
  headers := map[string]string{
    Authorization: fmt.Sprintf(tokenHeaderFormat, token),
    ContentType:   ApplicationJSON,
  }

  return client.requester(url, headers, body)
}

func resolveAPIError(body []byte, malformedError error, errorMap map[string]error) error {
  var errorMessage apiError
  apiErr := json.Unmarshal(body, &errorMessage)

  if apiErr != nil {
    return malformedError
  }

  err := errorMap[errorMessage.Error]

  if err == nil {
    return errors.ErrorNotIdentifiedFromAPI
  }

  return err
}
