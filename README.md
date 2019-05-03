# affiliate-marketing
Golang Client Library for Affiliate Integration

## Usage

Instantiate the client as below, 
being URL the Kueski API hosting, and the API/secret keys the credentials provided for each service.

```go
client := kueski.NewClient(url, apiKey, secretKey)
```

You should call client.Evaluate(curp, email, fullData), 
being curp and email a valid formatted string for each value, and fullData a struct containing
any extra data about the lead that you may want to pass to Kueski.
This data should be able to be marshaled to JSON, and in the future might help Kueski decide if
they want the lead or not.

```go
requestID, err := client.Evaluate(curp, email, fullData)

if (err != nil) {
  log.error(err.String())
}

doSomethingWithTheResponse(requestID, curp, email)
```

The `Evaluate` function will return the request ID (for tracking purposes) as string and an error code (if any).
Error codes can be resolved to string invoking the `String()` function.

Below you can find the list of available error codes.
Those marked with an * means that this case is handled successfully by the library.

| ErrorCode                           | ErrorNumber | Description |
|-------------------------------------|-------------|-------------|
| InvalidCurp                         | 1           | CURP format is invalid |
| InvalidEmail                        | 2           | Email format is invalid |
| InvalidCurpAndEmail                 | 3           | Both CURP and email formats are invalid |
| InvalidFullDataFormat               | 4           | Full data contains non JSON serializable data |
| InvalidRequestID                    | 5           | Request ID turned invalid * |
| InvalidFullDataAndRequestID         | 6           | Both full data and request ID are invalid |
| RequestIDNotFound                   | 7           | Request ID is not found in Kueski database |
| MissingCurp                         | 11          | CURP is missing * |
| MissingEmail                        | 12          | Email is missing * |
| MissingCurpAndEmail                 | 13          | Both CURP and Email are missing * |
| MissingCurpInvalidEmail             | 14          | CURP is missing and Email is invalid * |
| MissingEmailInvalidCurp             | 15          | Email is missing and CURP is invalid * |
| MissingFullData                     | 16          | Full data is nil |
| MissingRequestID                    | 17          | Request ID is missing when validated * |
| MissingFullDataAndRequestID         | 18          | Request ID and full data are missing * |
| MissingFullDataInvalidRequestID     | 19          | Full data is nil and Request ID is invalid * |
| MissingRequestIDInvalidFullData     | 20          | Request ID is missing and full data is non serializable * |
| UnableToMakeConnection              | 21          | Kueski host is unreachable |
| LeadEvaluationMalformedRequest      | 22          | Request for Lead Evaluation is malformed * |
| InvalidLeadEvaluationResponseFormat | 23          | Lead evaluation returned an error 500 |
| LeadDataMalformedRequest            | 24          | Request for Lead Data is malformed * |
| InvalidLeadDataResponseFormat       | 25          | Lead Data returned an error 500 |
| ErrorNotIdentifiedFromAPI           | 26          | Kueski API returned a non recognized validation error |
| AccessDenied                        | 31          | Any of API/Secret key are disabled or invalid |
| InvalidJWTResponseFormat            | 32          | Response from JWT request is malformed * |
| UnableToRefreshJWT                  | 33          | Kueski host is unreachable |
| InvalidSignatureFormat              | 34          | Authentication signature is malformed * |
| ExpiredJWTToken                     | 35          | Token is expired and unable to renew, a retry is needed |
| ExistingLead                        | 41          | Evaluated lead exists in Kueski database |
| DuplicatedLead                      | 42          | An evaluation with any of the CURP or email has been performed before |
| GeneralError                        | 99          | Generic error, it indicates an error in the library |
