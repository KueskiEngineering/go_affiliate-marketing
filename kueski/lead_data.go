package kueski

import (
  "encoding/json"

  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/errors"
  "github.com/KueskiEngineering/go_affiliate-marketing/kueski/util"
)

// leadDataPath - Path to the endpoint where to post lead full data.
const leadDataPath string = "affiliates/lead-data"

type leadFullData struct {
  Data      interface{} `json:"full_data"`
  RequestID string      `json:"request_id"`
}

type leadDataResponse struct {
  Response  string
  RequestID string `json:"request_id"`
}

var leadDataErrors = map[string]error{
  "Request not found":                           errors.RequestIDNotFound,
  "request_id is invalid, full_data is missing": errors.MissingFullDataInvalidRequestID,
  "request_id is invalid, full_data is invalid": errors.InvalidFullDataAndRequestID,
  "request_id is invalid":                       errors.InvalidRequestID,
  "request_id is missing, full_data is missing": errors.MissingFullDataAndRequestID,
  "request_id is missing, full_data is invalid": errors.MissingRequestIDInvalidFullData,
  "request_id is missing":                       errors.MissingRequestID,
  "full_data is missing":                        errors.MissingFullData,
  "full_data is invalid":                        errors.InvalidFullDataFormat,
}

// leadDataHandler - Type that defines the lead data handler signature.
type leadDataHandler func(client *Client, jsonData interface{}, requestID string) error

func leadData(client *Client, jsonData interface{}, requestID string) error {
  body, _ := json.Marshal(leadFullData{jsonData, requestID})
  response, err := client.makeRequest(leadDataPath, body)

  if err != nil {
    return err
  }

  responseBody, err := util.ExtractBody(response)

  if err != nil {
    return errors.InvalidLeadDataResponseFormat
  }

  if response.StatusCode == 500 {
    return errors.LeadDataMalformedRequest
  }

  if response.StatusCode == 400 {
    return resolveAPIError(responseBody, errors.InvalidLeadDataResponseFormat, leadDataErrors)
  }

  if response.StatusCode == 401 {
    return errors.ExpiredJWTToken
  }

  return resolveLeadDataResponse(responseBody, requestID)
}

func resolveLeadDataResponse(body []byte, requestID string) error {
  var response leadDataResponse
  responseErr := json.Unmarshal(body, &response)

  if responseErr != nil || !(response.RequestID == requestID && response.Response == "ok") {
    return errors.InvalidLeadDataResponseFormat
  }

  return nil
}
