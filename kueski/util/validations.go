package util

import (
  "encoding/json"
  "regexp"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
)

// ValidateEmail - Email validation.
func ValidateEmail(email string) bool {
  emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+" +
    "@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
    "[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
    "(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

  return emailRegex.MatchString(email)
}

// ValidateCurp - CURP validation.
func ValidateCurp(curp string) bool {
  curpRegex := regexp.MustCompile("^[A-Z]{4}\\d{2}(0[1-9]|10|11|12)" +
    "(0[1-9]|[1-2][0-9]|30|31)(H|M)(AS|BC|BS|CC|CS|CH|CL|CM|DF|DG" +
    "|MC|GT|GR|HG|JC|MN|MS|NT|NL|OC|PL|QT|QR|SP|SL|SR|TC|TS|TL|VZ" +
    "|YN|ZS|NE)[B-DF-HJ-NP-TV-Z]{3}[\\dA-Z]\\d$")
  return curpRegex.MatchString(curp)
}

// ValidateFullData - Validation of extra lead data.
func ValidateFullData(fullData interface{}) error {
  if fullData == nil {
    return errors.MissingFullData
  }

  _, err := json.Marshal(fullData)

  if err != nil {
    return errors.InvalidFullDataFormat
  }

  return nil
}
