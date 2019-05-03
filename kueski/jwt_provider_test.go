package kueski

import (
  "fmt"
  "math/rand"
  "testing"
  "time"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

var tokens = 10
var maxSeconds = 2
var tokenChannel = make(chan string, tokens)

// Mocks for different kind of testing.
type ConcurrentMock struct{}
type InvalidTokenClientMock struct{}
type InvalidRequesterClientMock struct{}
type SimpleClientMock struct{}

func (client *ConcurrentMock) tokenData() (string, int) {
  seconds := rand.Intn(maxSeconds)
  exp := int32(time.Now().Unix()) + int32(seconds)
  token := fmt.Sprintf(`{"token": "Valid Token", "expiration": %d}`, exp)
  return token, seconds
}

func (client *ConcurrentMock) RequestToken() ([]byte, error) {
  token, seconds := client.tokenData()
  time.Sleep(time.Duration(seconds) * time.Second)
  tokenChannel <- token
  return []byte(token), nil
}

func (client *InvalidTokenClientMock) RequestToken() ([]byte, error) {
  return []byte("Invalid Token"), nil
}

func (client *InvalidRequesterClientMock) RequestToken() ([]byte, error) {
  return nil, errors.GeneralError
}

func (client *SimpleClientMock) RequestToken() ([]byte, error) {
  return []byte(`{"token": "Token", "expiration": 12345678}`), nil
}

func TestTokenConcurrency(t *testing.T) {
  client := new(ConcurrentMock)
  provider := NewJWTProvider()

  for i := 0; i < tokens; i++ {
    go provider.Token(client)
  }

  for i := 0; i < tokens; i++ {
    token := <-tokenChannel
    assert.NotNil(t, token)
    assert.NotEqual(t, "", token)
  }
}

func TestTokenNotUpdated(t *testing.T) {
  client := new(ConcurrentMock)
  expiration := time.Now().AddDate(0, 0, 2)
  initialToken := fmt.Sprintf(`{"token": "Valid Token", "expiration": %d}`, int32(expiration.Unix()))

  provider := NewJWTProvider()
  provider.token = initialToken
  provider.exp = expiration

  for i := 0; i < tokens; i++ {
    token, err := provider.Token(client)
    assert.Nil(t, err)
    assert.NotNil(t, token)
    assert.Equal(t, initialToken, token)
  }
}

func TestInvalidTokenUnmarshalling(t *testing.T) {
  client := new(InvalidTokenClientMock)
  provider := NewJWTProvider()

  token, err := provider.Token(client)

  assert.Equal(t, "", token)
  assert.Equal(t, errors.InvalidJWTResponseFormat, err)
}

func TestInvalidTokenResponse(t *testing.T) {
  client := new(InvalidRequesterClientMock)
  provider := NewJWTProvider()

  token, err := provider.Token(client)

  assert.Equal(t, "", token)
  assert.Equal(t, errors.GeneralError, err)
}

func TestParseResponse(t *testing.T) {
  blob := []byte(`{"token": "Token", "expiration": 12345678}`)
  badBlob := []byte(`{"token": "Token", "expiration": "BadExpiration"}`)
  provider := NewJWTProvider()

  token, expiration, err := provider.parseResponse(blob)

  assert.Nil(t, err)
  assert.Equal(t, "Token", token)
  assert.Equal(t, 12345678, expiration)

  noToken, noExp, expErr := provider.parseResponse(badBlob)

  assert.Equal(t, "", noToken)
  assert.Equal(t, 0, noExp)
  assert.Equal(t, errors.InvalidJWTResponseFormat, expErr)
}
