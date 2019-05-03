package kueski

import (
  "encoding/json"
  "sync"
  "time"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
)

// TokenProvider - Interface that defines the Token retrieval signature.
type TokenProvider interface {
  Token(client TokenAccessor) (string, error)
}

// JWTProvider - Structure that handles the JWT requests.
type JWTProvider struct {
  token string
  exp   time.Time
  sync.Mutex
}

type authenticateResponse struct {
  Token      string
  Expiration int
}

var minutesToRenew = 10.0

// NewJWTProvider - JWT Provider constructor
func NewJWTProvider() *JWTProvider {
  jwt := new(JWTProvider)
  jwt.token = ""
  jwt.exp = time.Now()
  return jwt
}

// Token - Provides a valid JSON web Token.
func (jwt *JWTProvider) Token(client TokenAccessor) (string, error) {
  jwt.Lock()
  defer jwt.Unlock()

  // Check if its JWT renewal time
  if jwt.token == "" || time.Until(jwt.exp).Minutes() < minutesToRenew {
    blob, err := client.RequestToken()

    if err != nil {
      return "", err
    }

    token, expiration, err := jwt.parseResponse(blob)

    if err != nil {
      return "", err
    }

    jwt.token = token
    jwt.exp = time.Unix(int64(expiration), 0)
  }

  return jwt.token, nil
}

func (jwt *JWTProvider) parseResponse(blob []byte) (string, int, error) {
  var response authenticateResponse
  unmarshallErr := json.Unmarshal(blob, &response)

  if unmarshallErr != nil {
    return "", 0, errors.InvalidJWTResponseFormat
  }

  return response.Token, response.Expiration, nil
}
