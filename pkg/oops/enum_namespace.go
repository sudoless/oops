package oops

type Namespace uint8

const (
	NamespaceUnknown Namespace = iota

	// NamespaceInit assigned explicitly to errors arising from any init() function only
	NamespaceInit

	// NamespaceSetup assigned explicitly to errors occurring before the main program loop
	NamespaceSetup

	// NamespaceRuntime assigned to generic static functions that generate errors
	NamespaceRuntime

	// NamespaceIntegration assigned to third party integrations
	NamespaceIntegration

	// NamespaceApi assigned to errors occurring during the inbound API (e.g. http router)
	NamespaceApi

	// NamespaceIngress assigned to any other ingress source that may not necessarily be an api
	NamespaceIngress

	// NamespaceService assigned to core service logic occurring errors
	NamespaceService

	// NamespaceCache assigned to temporary/cache clients and logic layer (e.g. redis client)
	NamespaceCache

	// NamespaceStore assigned to persistent storage clients and logic layer (e.g. sql database client)
	NamespaceStore

	// NamespaceTest assigned only to Example, Test, Benchmark errors
	NamespaceTest

	namespaceMAX
)

func (e Namespace) String() string {
	code, ok := mapNamespaceToCode[e]
	if !ok {
		return "UNDEFINED"
	}
	return code
}

var mapNamespaceToCode = map[Namespace]string{
	NamespaceUnknown:     "UNKNOWN",
	NamespaceInit:        "INIT",
	NamespaceSetup:       "SETUP",
	NamespaceRuntime:     "RUNTIME",
	NamespaceIntegration: "INTEGRATION",
	NamespaceIngress:     "INGRESS",
	NamespaceApi:         "API",
	NamespaceService:     "SERVICE",
	NamespaceCache:       "CACHE",
	NamespaceStore:       "STORE",
	NamespaceTest:        "TEST",
}
