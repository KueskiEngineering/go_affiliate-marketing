package kueski

import (
  "bytes"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestAuthorizationToken(t *testing.T) {
  client := NewClient("", "apikey", "secretkey")
  assert.Equal(t, client.AuthorizationToken("canonical"), "apikey:1pbKbWCwwA/cOlxtE9+9L4wp4Bc=")
}

func TestErrorAtRequestToken(t *testing.T) {
  client := NewClient("URL", "Key", "Secret")
  errRequest := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    return nil, errors.GeneralError
  }
  client.requester = errRequest

  response, err := client.RequestToken()

  assert.Nil(t, response)
  assert.Equal(t, errors.UnableToRefreshJWT, err)
}

func TestMalformedHeaderAtRequestToken(t *testing.T) {
  client := NewClient("URL", "Key", "Secret")
  errRequest := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    response := http.Response{}
    response.StatusCode = 400
    return &response, nil
  }
  client.requester = errRequest

  response, err := client.RequestToken()

  assert.Nil(t, response)
  assert.Equal(t, errors.InvalidSignatureFormat, err)
}

func TestAccessDeniedAtRequestToken(t *testing.T) {
  client := NewClient("URL", "Key", "Secret")
  errRequest := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    response := http.Response{}
    response.StatusCode = 401
    return &response, nil
  }
  client.requester = errRequest

  response, err := client.RequestToken()

  assert.Nil(t, response)
  assert.Equal(t, errors.AccessDenied, err)
}

func TestErrorUnmarshalingAtRequestToken(t *testing.T) {
  client := NewClient("URL", "Key", "Secret")
  errRequest := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    handler := func(writer http.ResponseWriter, request *http.Request) {
      writer.Header().Set("Content-Length", "1")
    }

    testServer := httptest.NewServer(http.HandlerFunc(handler))
    defer testServer.Close()

    request, _ := http.NewRequest("GET", testServer.URL, bytes.NewBuffer([]byte("")))

    client := &http.Client{}
    response, _ := client.Do(request)

    return response, nil
  }

  client.requester = errRequest
  response, err := client.RequestToken()

  assert.Nil(t, response)
  assert.Equal(t, errors.InvalidJWTResponseFormat, err)
}

func TestSuccessfulRequestToken(t *testing.T) {
  client := NewClient("URL", "Key", "Secret")
  errRequest := func(url string, headers map[string]string, body []byte) (*http.Response, error) {
    response := http.Response{}
    response.StatusCode = 201
    response.Body = ioutil.NopCloser(strings.NewReader("Body"))
    return &response, nil
  }
  client.requester = errRequest

  response, err := client.RequestToken()

  assert.Nil(t, err)
  assert.Equal(t, "Body", string(response[:]))
}
