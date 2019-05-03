package util

import (
  "crypto/md5"
  "encoding/base64"
  "strings"
  "time"
)

// gmtLocation - GMT time zone.
var gmtLocation, _ = time.LoadLocation("GMT")

// Canonical - Builds the canonical string required for the signature.
func Canonical(method, contentType, content, urlPath string, now time.Time) string {
  methodUpper := strings.ToUpper(method)
  dateFormat := HTTPDate(now)
  contentEncoded := ContentMD5(content)
  return methodUpper + "," + contentType + "," + contentEncoded + "," + urlPath + "," + dateFormat
}

// ContentMD5 - Generate the MD5 of the Body string.
func ContentMD5(content string) string {
  contentData := []byte(content)
  md5Content := md5.Sum(contentData)
  var contentBytes = md5Content[:]
  contentEncoded := base64.StdEncoding.EncodeToString(contentBytes)
  return contentEncoded
}

// HTTPDate - Generate the Date in HTTP format for UTC/GMT time.
func HTTPDate(date time.Time) string {
  return date.In(gmtLocation).Format(time.RFC1123)
}
