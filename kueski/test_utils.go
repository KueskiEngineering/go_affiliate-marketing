package kueski

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
)

type fakeTokenProvider struct {
	works bool
}

func (provider *fakeTokenProvider) Token(client TokenAccessor) (string, error) {
	if provider.works {
		return "Token", nil
	}

	return "", errors.GeneralError
}

type invalidReader int

func (invalidReader) Read(p []byte) (_ int, err error) {
	return 0, errors.GeneralError
}

func invalidHTTPResponse() *http.Response {
	invalidBody := ioutil.NopCloser(invalidReader(0))

	response := http.Response{}
	response.Body = invalidBody

	return &response
}

func buildHTTPResponse(statusCode int, body string) *http.Response {
	response := http.Response{}
	response.StatusCode = statusCode
	response.Body = ioutil.NopCloser(strings.NewReader(body))
	return &response
}
