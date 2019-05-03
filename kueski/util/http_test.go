package util

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
  assert.Equal(t, "example.com/path", BuildURL("example.com", "path"))
}

type jsonTestMessage struct {
  StringVar string
  IntVar    int
}

func TestHttpPostRequest(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
    assert.Equal(t, request.Header.Get("Hello"), "World")
    assert.Equal(t, request.Header.Get("Foo"), "Bar")

    bodyBytes, _ := ioutil.ReadAll(request.Body)
    bodyString := string(bodyBytes)
    assert.Equal(t, bodyString, `{"StringVar":"Hello","IntVar":123456}`)

    //Response
    responseWriter.Header().Set("res-Hello", "res-World")
    responseWriter.Header().Set("res-Foo", "res-Bar")
    responseWriter.Write([]byte("This is the response body"))
  }))
  defer ts.Close()

  //Request
  msg := jsonTestMessage{"Hello", 123456}
  bodyMsg, _ := json.Marshal(msg)

  headers := make(map[string]string)
  headers["Hello"] = "World"
  headers["Foo"] = "Bar"

  res, _ := PostRequest(ts.URL, headers, bodyMsg)

  assert.Equal(t, res.Header.Get("res-Hello"), "res-World")
  assert.Equal(t, res.Header.Get("res-Foo"), "res-Bar")

  defer res.Body.Close()
  bodyBytes, _ := ioutil.ReadAll(res.Body)
  bodyString := string(bodyBytes)
  assert.Equal(t, bodyString, "This is the response body")
}

func TestErrorAtExtractBody(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
    //Response
    responseWriter.Header().Set("Content-Length", "1")
    responseWriter.Write([]byte(""))
  }))
  defer ts.Close()

  body, _ := json.Marshal("")
  headers := make(map[string]string)
  response, _ := PostRequest(ts.URL, headers, body)

  body, err := ExtractBody(response)

  assert.Nil(t, body)
  assert.NotNil(t, err)
}

func TestExtractBody(t *testing.T) {
  expected := "Body"
  response := http.Response{}
  response.Body = ioutil.NopCloser(strings.NewReader(expected))

  body, err := ExtractBody(&response)

  assert.Nil(t, err)
  assert.Equal(t, []byte(expected), body)
}

func TestErroneousPostRequest(t *testing.T) {
  body, _ := json.Marshal("")
  headers := make(map[string]string)
  response, err := PostRequest("300.400.500.600", headers, body)

  assert.Nil(t, response)
  assert.Equal(t, err, errors.UnableToMakeConnection)
}
