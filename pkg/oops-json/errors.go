package oops_json

import "go.sdls.io/oops/pkg/oops"

var (
	ErrInvalid = oops.Define(oops.BlameClient, oops.NamespaceRuntime, oops.ReasonResourceDecoding,
		"invalid json syntax, please use a validator to check the json syntax")
	ErrDecoding = oops.Define(oops.BlameClient, oops.NamespaceRuntime, oops.ReasonResourceDecoding,
		"failed to decode json, please ensure you're using the right types")
	ErrEncoding = oops.Define(oops.BlameClient, oops.NamespaceRuntime, oops.ReasonResourceEncoding,
		"failed to encode json, please ensure you're using the right types")
	ErrUnexpected = oops.Define(oops.BlameDeveloper, oops.NamespaceRuntime, oops.ReasonUnexpected)
)
