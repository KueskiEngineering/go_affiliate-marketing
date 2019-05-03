package errors

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestResponseError(t *testing.T) {
  for responseError, description := range descriptions {
    assert.Equal(t, description.name, responseError.String())
    assert.Equal(t, description.name, responseError.Error())
    assert.Equal(t, description.desc, responseError.Description())
    assert.Equal(t, responseError, StringToResponseError(description.name))
  }

  unknown := ResponseError(10)
  assert.Equal(t, "ResponseError(10)", unknown.String())
  assert.Equal(t, "ResponseError(10)", unknown.Error())
  assert.Equal(t, "Unknown response error.", unknown.Description())
  assert.Nil(t, StringToResponseError("This is not a valid error"))
}
