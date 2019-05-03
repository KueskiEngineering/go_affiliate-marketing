package kueski

import (
  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
)

type curpValidator func(curp string) bool
type emailValidator func(email string) bool
type fullDataValidator func(fullData interface{}) error

// Validator - Affiliate data validator: CURP, email and full data.
type Validator struct {
  curp  curpValidator
  email emailValidator
  data  fullDataValidator
}

func (validator *Validator) validate(curp, email string, fullData interface{}) error {
  curpValid := validator.curp(curp)
  emailValid := validator.email(email)

  if !curpValid && !emailValid {
    return errors.InvalidCurpAndEmail
  }

  if !curpValid {
    return errors.InvalidCurp
  }

  if !emailValid {
    return errors.InvalidEmail
  }

  return validator.data(fullData)
}
