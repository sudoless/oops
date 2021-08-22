package oops

import (
	"fmt"
	"net/http"
	"testing"
)

func TestReason_enumMapping(t *testing.T) {
	t.Parallel()

	var reason Reason

	for reason = ReasonUnknown; reason < reasonMAX; reason++ {
		t.Run(fmt.Sprintf("reason(%d)", reason), func(t *testing.T) {
			if reason.String() == "" {
				t.Error("reason cannot have empty String()")
				t.Logf("")
			}

			if code := reason.HttpStatusCode(); code == 0 || code == http.StatusTeapot {
				t.Error("reason cannot have 0 value HttpStatusCode()")
			}

			if t.Failed() {
				t.Logf("after:  '%s'", (reason - 1).String())
				t.Logf("before: '%s'", (reason + 1).String())
			}
		})
	}
}

func TestBlame_enumMapping(t *testing.T) {
	t.Parallel()

	var blame Blame

	for blame = BlameUnknown; blame < blameMAX; blame++ {
		t.Run(fmt.Sprintf("blame(%d)", blame), func(t *testing.T) {
			if blame.String() == "" {
				t.Error("blame cannot have empty String()")
				t.Logf("")
			}

			if t.Failed() {
				t.Logf("after:  '%s'", (blame - 1).String())
				t.Logf("before: '%s'", (blame + 1).String())
			}
		})
	}
}

func TestNamespace_enumMapping(t *testing.T) {
	t.Parallel()

	var namespace Namespace

	for namespace = NamespaceUnknown; namespace < namespaceMAX; namespace++ {
		t.Run(fmt.Sprintf("namespace(%d)", namespace), func(t *testing.T) {
			if namespace.String() == "" {
				t.Error("namespace cannot have empty String()")
				t.Logf("")
			}

			if t.Failed() {
				t.Logf("after:  '%s'", (namespace - 1).String())
				t.Logf("before: '%s'", (namespace + 1).String())
			}
		})
	}
}

func Test_enumNotMapped(t *testing.T) {
	if (reasonMAX + 1).String() != "UNDEFINED" {
		t.Error("reasonMAX+1 must be UNDEFINED")
	}
	if (reasonMAX + 1).HttpStatusCode() != http.StatusTeapot {
		t.Error("reasonMAX+1 must be 418 http status code")
	}
	if (blameMAX + 1).String() != "UNDEFINED" {
		t.Error("blameMAX+1 must be UNDEFINED")
	}
	if (namespaceMAX + 1).String() != "UNDEFINED" {
		t.Error("namespaceMAX+1 must be UNDEFINED")
	}
}
