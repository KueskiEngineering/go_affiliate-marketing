package util

import (
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
)

func TestCanonical(t *testing.T) {
  str := "Fri, 25 Jan 2019 18:21:25 GMT"
  now, _ := time.Parse(time.RFC1123, str)
  assert.Equal(t, Canonical("posT", "type", "content", "url", now),
    "POST,type,mgNkuembtIDdJeHwKEyFVQ==,url,Fri, 25 Jan 2019 18:21:25 GMT")
}

func TestContentMD5(t *testing.T) {
  assert.Equal(t, "1B2M2Y8AsgTpgAmY7PhCfg==", ContentMD5(""))
}

func TestHTTPDate(t *testing.T) {
  str := "Fri, 25 Jan 2019 18:21:25 GMT"
  now, _ := time.Parse(time.RFC1123, str)
  assert.Equal(t, str, HTTPDate(now))
}
