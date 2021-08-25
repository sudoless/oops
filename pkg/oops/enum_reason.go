package oops

import "net/http"

type Reason uint8

const (
	// ReasonUnknown MUST NOT be used in code, it acts as a way to detect badly decoded, encoded, etc errors
	ReasonUnknown Reason = iota

	// ReasonUnexpected for unexpected failures from an unknown point or reason (e.g. panic)
	ReasonUnexpected

	// ReasonInternal for unexpected failures, but from a known point (e.g. detecting a nil pointer)
	ReasonInternal

	// ReasonUnavailable for when the server is not healthy and service degradation is at play
	ReasonUnavailable

	// ReasonConnection for any sort of failed connections (raw TCP, database or cache connection, internet, etc)
	ReasonConnection

	// ReasonTimeout for any connection or context timeout error
	ReasonTimeout

	// ReasonIO for any io.Writer... or io.Reader... interaction that failed
	ReasonIO

	// ReasonOS for any operating system call, file system interaction, etc
	ReasonOS

	// ReasonValidation for failing github.com/go-playground/validator validation
	ReasonValidation

	// ReasonValidationLookup for field that may pass validation rules but do not pass a lookup/search (e.g. city name)
	ReasonValidationLookup

	// ReasonDbGeneric for any database error that does not fit into the other reasons
	ReasonDbGeneric

	// ReasonDbQuery for failing on Query... function
	ReasonDbQuery

	// ReasonDbExec for failing on Exec... function
	ReasonDbExec

	// ReasonDbScan for failing on Scan... function
	ReasonDbScan

	// ReasonDbTx for failing on initializing a database transaction
	ReasonDbTx

	// ReasonDbTxCommit for failing to commit a database transaction
	ReasonDbTxCommit

	// ReasonDbTxRollback for failing to rollback a database transaction
	ReasonDbTxRollback

	// ReasonAuthNone for missing authentication/authorization credentials
	ReasonAuthNone

	// ReasonAuthFormat for providing authentication/authorization credentials that are invalid or badly formatted
	ReasonAuthFormat

	// ReasonAuthBad for providing authentication credentials that are not valid
	ReasonAuthBad

	// ReasonAuthForbidden for requesting a resource with authorization credentials that not allowed
	ReasonAuthForbidden

	// ReasonRateLimit for warning that the request was rate limited and the client should try again later
	ReasonRateLimit

	// ReasonRateLimitBan for warning that the client has been banned from future requests
	ReasonRateLimitBan

	// ReasonLegal for rejecting requests that do not meet legal requirements (e.g. banning UK requests due to Brexit)
	ReasonLegal

	// ReasonResourceEncoding for any resource encoding time errors (JSON, XML, BASE64, HEX, keys, etc)
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceEncoding

	// ReasonResourceDecoding for any resource decoding time errors (JSON, XML, BASE64, HEX, keys, etc)
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceDecoding

	// ReasonResourceNotFound for when a resource cannot be found
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceNotFound

	// ReasonResourceGone for when a resource cannot be found but is known to have existed or for when a resource
	// was found but it is no longer "active" or it has "expired"
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceGone

	// ReasonResourceNotYet for when a resource cannot be found at the moment, but it is possible that it will be
	// available in the future
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceNotYet

	// ReasonResourceTooLarge for when a resource cannot be processed due to being too large
	// (Resource in this context refers to data/info/object/etc being returned BY the server TO the client)
	ReasonResourceTooLarge

	// ReasonPayloadEncoding for any payload encoding time errors (JSON, XML, BASE64, HEX, keys, etc)
	// (Payload in this context refers to data/info/object/etc being sent BY the server TO another server/service)
	ReasonPayloadEncoding

	// ReasonPayloadDecoding for any payload decoding time errors (JSON, XML, BASE64, HEX, keys, etc)
	// (Payload in this context refers to data/info/object/etc being sent BY the server TO another server/service)
	ReasonPayloadDecoding

	// ReasonPayloadTooLarge for when a payload cannot be processed due to being too large
	// (Payload in this context refers to data/info/object/etc being sent BY the server TO another server/service)
	ReasonPayloadTooLarge

	// ReasonRequestFormat for when the request does not meet the expected format or it is improperly encoded
	ReasonRequestFormat

	// ReasonRequestDecoding for when a request body/data cannot be decoded
	ReasonRequestDecoding

	// ReasonRequestTooLarge for when the request size is too large and cannot and will not be processed
	ReasonRequestTooLarge

	// ReasonRequestMissing for when there is request information missing
	ReasonRequestMissing

	// ReasonRequestBad for generic bad request, should avoid using this reason
	ReasonRequestBad

	// ReasonRequestConflict for when a request was understood, but the action cannot be performed due to constraints
	ReasonRequestConflict

	// ReasonRequestUnprocessable for when a request was understood (format and protocol wise) but the instructions
	// were not understood
	ReasonRequestUnprocessable

	// ReasonRequestValidationParameters for when validating a requests parameters (path, path args, query args, headers,
	// body, etc) outside of using the Validator and ReasonValidation
	ReasonRequestValidationParameters

	// ReasonRequestMethodNotAllowed for when a method is not allowed and a different method should be tried
	ReasonRequestMethodNotAllowed

	// ReasonRequestEndpointNotFound for when an endpoint cannot be found by the given path
	ReasonRequestEndpointNotFound

	// ReasonIdempotency for reporting idempotency inconsistencies (race conditions, key reuse on different endpoints)
	ReasonIdempotency

	// ReasonConfig for when a configuration or defined (or pre-defined) state is not valid, OR they do not match
	// expectation or requested actions
	ReasonConfig

	// ReasonConfigMissing for when a configuration or expected state is not present
	ReasonConfigMissing

	// ReasonCrypto for any generic crypto issues (bad key, invalid size, missing nonce, etc)
	ReasonCrypto

	// ReasonGatewayUnavailable for when the gateway is not available or able to process a request
	ReasonGatewayUnavailable

	// ReasonGatewayForwarding for when the gateway was unable to forward (in any direction) a request/response
	ReasonGatewayForwarding

	// ReasonGatewayAuth for when a given request is not authorized to access the gateway
	ReasonGatewayAuth

	// ReasonGatewayFailure for internal gateway errors
	ReasonGatewayFailure

	// ReasonCORS for any API CORS errors
	ReasonCORS

	// reasonMAX acts as an internal testing landmark, to check that all enums before it have the necessary map value
	reasonMAX
)

func (e Reason) String() string {
	code, ok := mapReasonToCode[e]
	if !ok {
		return "UNDEFINED"
	}
	return code
}

var mapReasonToCode = map[Reason]string{
	ReasonUnknown: "UNKNOWN",

	ReasonUnexpected:  "UNEXPECTED",
	ReasonInternal:    "INTERNAL",
	ReasonUnavailable: "UNAVAILABLE",
	ReasonConnection:  "CONNECTION",
	ReasonTimeout:     "TIMEOUT",
	ReasonIO:          "IO",
	ReasonOS:          "OS",

	ReasonValidation:       "VALIDATION",
	ReasonValidationLookup: "VALIDATION_LOOKUP",

	ReasonDbGeneric:    "DB_GENERIC",
	ReasonDbQuery:      "DB_QUERY",
	ReasonDbExec:       "DB_EXEC",
	ReasonDbScan:       "DB_SCAN",
	ReasonDbTx:         "DB_TX",
	ReasonDbTxCommit:   "DB_TX_COMMIT",
	ReasonDbTxRollback: "DB_TX_ROLLBACK",

	ReasonAuthNone:      "AUTH_NONE",
	ReasonAuthFormat:    "AUTH_FORMAT",
	ReasonAuthBad:       "AUTH_BAD",
	ReasonAuthForbidden: "AUTH_FORBIDDEN",

	ReasonRateLimit:    "RATE_LIMIT",
	ReasonRateLimitBan: "RATE_LIMIT_BAN",

	ReasonLegal: "LEGAL",

	ReasonResourceEncoding: "RESOURCE_ENCODING",
	ReasonResourceDecoding: "RESOURCE_DECODING",
	ReasonResourceNotFound: "RESOURCE_NOT_FOUND",
	ReasonResourceGone:     "RESOURCE_GONE",
	ReasonResourceNotYet:   "RESOURCE_NOT_YET",
	ReasonResourceTooLarge: "RESOURCE_TOO_LARGE",

	ReasonPayloadEncoding: "PAYLOAD_ENCODING",
	ReasonPayloadDecoding: "PAYLOAD_DECODING",
	ReasonPayloadTooLarge: "PAYLOAD_TOO_LARGE",

	ReasonRequestFormat:        "REQUEST_FORMAT",
	ReasonRequestTooLarge:      "REQUEST_TOO_LARGE",
	ReasonRequestMissing:       "REQUEST_MISSING",
	ReasonRequestDecoding:      "REQUEST_DECODING",
	ReasonRequestBad:           "REQUEST_BAD",
	ReasonRequestConflict:      "REQUEST_CONFLICT",
	ReasonRequestUnprocessable: "REQUEST_UNPROCESSABLE",

	ReasonRequestValidationParameters: "REQUEST_VALIDATION_PARAMETERS",
	ReasonRequestMethodNotAllowed:     "REQUEST_METHOD_NOT_ALLOWED",
	ReasonRequestEndpointNotFound:     "REQUEST_ENDPOINT_NOT_FOUND",

	ReasonIdempotency: "IDEMPOTENCY",

	ReasonConfig:        "CONFIG",
	ReasonConfigMissing: "CONFIG_MISSING",

	ReasonCrypto: "CRYPTO",

	ReasonGatewayUnavailable: "GATEWAY_UNAVAILABLE",
	ReasonGatewayForwarding:  "GATEWAY_FORWARDING",
	ReasonGatewayAuth:        "GATEWAY_AUTH",
	ReasonGatewayFailure:     "GATEWAY_FAILURE",

	ReasonCORS: "CORS",
}

func (e Reason) HttpStatusCode() int {
	code, ok := mapReasonToHttpStatus[e]
	if !ok {
		return http.StatusTeapot
	}
	return code
}

var mapReasonToHttpStatus = map[Reason]int{
	ReasonUnknown: http.StatusInternalServerError,

	ReasonUnexpected:  http.StatusInternalServerError,
	ReasonInternal:    http.StatusInternalServerError,
	ReasonUnavailable: http.StatusServiceUnavailable,
	ReasonConnection:  http.StatusInternalServerError,
	ReasonTimeout:     http.StatusGatewayTimeout,
	ReasonIO:          http.StatusInternalServerError,
	ReasonOS:          http.StatusInternalServerError,

	ReasonValidation:       http.StatusBadRequest,
	ReasonValidationLookup: http.StatusBadRequest,

	ReasonDbGeneric:    http.StatusInternalServerError,
	ReasonDbQuery:      http.StatusInternalServerError,
	ReasonDbExec:       http.StatusInternalServerError,
	ReasonDbScan:       http.StatusInternalServerError,
	ReasonDbTx:         http.StatusInternalServerError,
	ReasonDbTxCommit:   http.StatusInternalServerError,
	ReasonDbTxRollback: http.StatusInternalServerError,

	ReasonAuthNone:      http.StatusUnauthorized,
	ReasonAuthFormat:    http.StatusUnauthorized,
	ReasonAuthBad:       http.StatusUnauthorized,
	ReasonAuthForbidden: http.StatusForbidden,

	ReasonRateLimit:    http.StatusTooManyRequests,
	ReasonRateLimitBan: http.StatusTooManyRequests,

	ReasonLegal: http.StatusUnavailableForLegalReasons,

	ReasonResourceEncoding: http.StatusInternalServerError,
	ReasonResourceDecoding: http.StatusInternalServerError,
	ReasonResourceNotFound: http.StatusNotFound,
	ReasonResourceGone:     http.StatusGone,
	ReasonResourceNotYet:   http.StatusTooEarly,
	ReasonResourceTooLarge: http.StatusRequestEntityTooLarge,

	ReasonPayloadEncoding: http.StatusInternalServerError,
	ReasonPayloadDecoding: http.StatusInternalServerError,
	ReasonPayloadTooLarge: http.StatusInternalServerError,

	ReasonRequestFormat:        http.StatusUnsupportedMediaType,
	ReasonRequestTooLarge:      http.StatusRequestEntityTooLarge,
	ReasonRequestMissing:       http.StatusMisdirectedRequest,
	ReasonRequestDecoding:      http.StatusBadRequest,
	ReasonRequestBad:           http.StatusBadRequest,
	ReasonRequestConflict:      http.StatusConflict,
	ReasonRequestUnprocessable: http.StatusUnprocessableEntity,

	ReasonRequestValidationParameters: http.StatusBadRequest,
	ReasonRequestMethodNotAllowed:     http.StatusMethodNotAllowed,
	ReasonRequestEndpointNotFound:     http.StatusNotFound,

	ReasonIdempotency: http.StatusLocked,

	ReasonConfig:        http.StatusConflict,
	ReasonConfigMissing: http.StatusBadRequest,

	ReasonCrypto: http.StatusBadRequest,

	ReasonGatewayUnavailable: http.StatusServiceUnavailable,
	ReasonGatewayForwarding:  http.StatusBadGateway,
	ReasonGatewayAuth:        http.StatusProxyAuthRequired,
	ReasonGatewayFailure:     http.StatusBadGateway,

	ReasonCORS: http.StatusForbidden,
}

var mapCodeToReason = func() map[string]Reason {
	ret := make(map[string]Reason, len(mapReasonToCode))
	for k, v := range mapReasonToCode {
		ret[v] = k
	}
	return ret
}()
