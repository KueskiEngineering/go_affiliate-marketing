package util

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "net/http"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
)

// PostRequestFunc - Interface for web requests.
type PostRequestFunc func(url string, headers map[string]string, body []byte) (*http.Response, error)

// BuildURL - URL builder.
func BuildURL(host, path string) string {
  return fmt.Sprintf("%s/%s", host, path)
}

// PostRequest - HTTP POST Request using Golang stdlib.
// url - String URL.
// headers - Headers to be included in the request.
// body - Request body.
func PostRequest(url string, headers map[string]string, body []byte) (*http.Response, error) {
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

  for key, value := range headers {
    req.Header.Set(key, value)
  }

  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    return nil, errors.UnableToMakeConnection
  }

  return resp, err
}

// ExtractBody - Encapsulation of the task which extracts body response data.
func ExtractBody(response *http.Response) ([]byte, error) {
  defer response.Body.Close()
  responseBody, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return nil, err
  }

  return responseBody, nil
}
