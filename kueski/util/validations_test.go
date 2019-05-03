package util

import (
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
  assert.False(t, ValidateEmail(""))
  assert.True(t, ValidateEmail("test@kueski.com"))
  assert.True(t, ValidateEmail("test@kueski-one.com"))
  assert.False(t, ValidateEmail("test@kueski_one.com"))
}

func TestCurpValidation(t *testing.T) {
  assert.True(t, ValidateCurp("ABCD920113MSLXYZ01"))
  assert.False(t, ValidateCurp("ABCD920113MSLXY101"))
  assert.False(t, ValidateCurp("AB2011CD93MSLXY101"))
  assert.False(t, ValidateCurp("AB201SLXY101"))
}

type AbleToBeConvertedJSON struct {
  A string
  b string
  C int
}

type UnableToBeConvertedJSON struct {
  A func(a string) string
  B string
}

func TestFullDataValidation(t *testing.T) {
  function := func(a string) string { return "S" }
  able := AbleToBeConvertedJSON{"A", "b", 1}
  unable := UnableToBeConvertedJSON{function, "B"}

  noErr := ValidateFullData(able)
  err := ValidateFullData(unable)
  nilErr := ValidateFullData(nil)

  assert.Nil(t, noErr)
  assert.Equal(t, errors.MissingFullData, nilErr)
  assert.Equal(t, errors.InvalidFullDataFormat, err)
}
