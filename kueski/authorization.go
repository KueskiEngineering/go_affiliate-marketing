package kueski

import (
  "crypto/hmac"
  "crypto/sha1"
  "encoding/base64"
  "fmt"
  "time"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/util"
)

// ContentMD5       - Content MD5 header key.
// Authorization    - Authorization header key.
// Date             - Date header key.
// ContentType      - Content Type header key.
// ApplicationJSON  - application/JSON string value.
// AuthenticatePath - Path to the authentication signature endpoint.
// BodyString       - Authentication request body.
// Method           - HTTP Method to use at Authentication (POST).
const (
  ContentMD5       string = "Content-MD5"
  Authorization    string = "Authorization"
  Date             string = "Date"
  ContentType      string = "Content-Type"
  ApplicationJSON  string = "application/json"
  AuthenticatePath string = "affiliates/authenticate"
  BodyString       string = ""
  Method           string = "POST"
  APIAuthPrefix    string = "APIAuth"
)

// TokenAccessor - Interface that defines the API auth methods.
type TokenAccessor interface {
  RequestToken() ([]byte, error)
}

// AuthorizationToken - Creates a signature for API authentication.
func (client *Client) AuthorizationToken(canonical string) string {
  keyForSign := []byte(client.secretKey)
  h := hmac.New(sha1.New, keyForSign)
  h.Write([]byte(canonical))
  return client.apiKey + ":" + base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// RequestToken - Internal function to retrieve a valid JWT from Kueski API
func (client *Client) RequestToken() ([]byte, error) {
  url := util.BuildURL(client.url, AuthenticatePath)
  body := []byte(BodyString)
  now := time.Now()
  httpDate := util.HTTPDate(now)

  canonical := util.Canonical(Method, ApplicationJSON, BodyString, fmt.Sprintf("/%s", AuthenticatePath), now)

  headers := map[string]string{
    ContentMD5:    util.ContentMD5(BodyString),
    Authorization: fmt.Sprintf("%s %s", APIAuthPrefix, client.AuthorizationToken(canonical)),
    Date:          httpDate,
    ContentType:   ApplicationJSON,
  }

  response, err := client.requester(url, headers, body)

  if err != nil {
    return nil, errors.UnableToRefreshJWT
  }

  if response.StatusCode == 400 {
    return nil, errors.InvalidSignatureFormat
  }

  if response.StatusCode == 401 {
    return nil, errors.AccessDenied
  }

  responseBody, err := util.ExtractBody(response)

  if err != nil {
    return nil, errors.InvalidJWTResponseFormat
  }

  return responseBody, nil
}
