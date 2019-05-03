package kueski

import (
  "testing"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
  failedString := func(s string) bool { return false }
  failedData := func(s interface{}) error { return errors.GeneralError }
  goodString := func(s string) bool { return true }
  goodData := func(s interface{}) error { return nil }

  validators := []Validator{
    Validator{failedString, failedString, goodData},
    Validator{failedString, goodString, goodData},
    Validator{goodString, failedString, goodData},
    Validator{goodString, goodString, failedData},
    Validator{goodString, goodString, goodData},
  }
  results := []error{errors.InvalidCurpAndEmail, errors.InvalidCurp, errors.InvalidEmail, errors.GeneralError, nil}

  for i, validator := range validators {
    assert.Equal(t, results[i], validator.validate("", "", validator))
  }
}
