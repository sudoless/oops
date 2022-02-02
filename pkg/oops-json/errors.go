package oops_json

import (
	"net/http"

	"go.sdls.io/oops/pkg/oops"
)

var (
	errGroup    = oops.Define().StatusCode(http.StatusBadRequest).Type("json").Group().PrefixCode("json_")
	ErrInvalid  = errGroup.Code("invalid").Help("invalid json syntax, please use a validator to check the json syntax")
	ErrDecoding = errGroup.Code("decode").Help("failed to decode json, please ensure you're using the right types")
	ErrEncoding = errGroup.Code("encode").Help("failed to encode json, please ensure you're using the right types")
)
