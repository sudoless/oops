package oops_test

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestError_Format(t *testing.T) {
	t.Parallel()

	def := oops.Define("load").Message("load failed").Causes(oops.CauseIO).Actions(oops.ActionRetry)
	child := oops.Define("item").Causes(oops.CauseNotFound).Yeet().Pathf("items/%d", 3)

	t.Run("v s and q print Error", func(t *testing.T) {
		t.Parallel()
		err := def.Yeetf("config").Nest(errors.New("hidden"))
		for _, tc := range []struct{ format, want string }{
			{"%v", "load: load failed; config"},
			{"%s", "load: load failed; config"},
			{"%q", `"load: load failed; config"`},
			{"[%30s]", "[     load: load failed; config]"},
		} {
			if got := fmt.Sprintf(tc.format, err); got != tc.want {
				t.Errorf("%s: got %q, want %q", tc.format, got, tc.want)
			}
		}
	})

	t.Run("+v prints the tree", func(t *testing.T) {
		t.Parallel()
		err := def.Yeetf("config").Pathf("root").Set("b", 2).Set("a", "x").
			Nest(child.Nest(fmt.Errorf("open: %w", io.EOF)), io.ErrClosedPipe)
		want := strings.Join([]string{
			"load: load failed; config",
			"  path: root",
			"  causes: io",
			"  actions: retry",
			"  fields: a=x, b=2",
			"  └ item",
			"    path: items/3",
			"    causes: not_found",
			"    └ open: EOF",
			"  └ io: read/write on closed pipe",
		}, "\n")
		if got := fmt.Sprintf("%+v", err); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("+v prints a foreign subtree once, indenting multi-line text", func(t *testing.T) {
		t.Parallel()
		join := errors.Join(io.EOF, oops.Define("inside").Yeet())
		after := oops.Define("after").Yeet().Nest(oops.Define("below").Yeet())
		err := oops.Define("outer").Wrap(join).Nest(after)
		want := strings.Join([]string{
			"outer",
			"  └ EOF",
			"    inside",
			"  └ after",
			"    └ below",
		}, "\n")
		if got := fmt.Sprintf("%+v", err); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("+v indents multi-line field values", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("outer").Yeet().Nest(oops.Define("inner").Yeet().Set("query", "SELECT 1\nFROM t"))
		want := strings.Join([]string{
			"outer",
			"  └ inner",
			"    fields: query=SELECT 1",
			"    FROM t",
		}, "\n")
		if got := fmt.Sprintf("%+v", err); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("+v ends with a line when the tree exceeds the node cap", func(t *testing.T) {
		t.Parallel()
		fan := func(children int) *oops.Error {
			errs := make([]error, children)
			for i := range errs {
				errs[i] = oops.Define("leaf").Yeet()
			}
			return oops.Define("root").Yeet().Nest(errs...)
		}

		lines := strings.Split(fmt.Sprintf("%+v", fan(1023)), "\n")
		if len(lines) != 1024 || lines[1023] != "  └ leaf" {
			t.Fatalf("1024 nodes: got %d lines ending %q", len(lines), lines[len(lines)-1])
		}

		lines = strings.Split(fmt.Sprintf("%+v", fan(1024)), "\n")
		if len(lines) != 1025 || lines[1023] != "  └ leaf" || lines[1024] != "  … truncated after 1024 nodes" {
			t.Fatalf("1025 nodes: got %d lines ending %q", len(lines), lines[len(lines)-2:])
		}
	})

	t.Run("+v uses the definition formatter for the header", func(t *testing.T) {
		t.Parallel()
		custom := oops.Define("custom").Formatter(func(e *oops.Error) string { return "custom<" + e.Code() + ">" })
		if got := fmt.Sprintf("%+v", custom.Yeet()); got != "custom<custom>" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("+v lists trace frames", func(t *testing.T) {
		t.Parallel()
		err := errTraced.Yeet()
		lines := strings.Split(fmt.Sprintf("%+v", err), "\n")
		if len(lines) < 3 || lines[0] != "trace.test" || lines[1] != "  trace:" {
			t.Fatalf("got %q", lines)
		}
		if len(lines)-2 != len(err.Trace()) {
			t.Fatalf("got %d frame lines, want %d", len(lines)-2, len(err.Trace()))
		}
		for i, line := range lines[2:] {
			if line != "    "+err.Trace()[i] {
				t.Fatalf("line %d = %q, want frame %q", i+2, line, err.Trace()[i])
			}
		}
		if !regexp.MustCompile(`: TestError_Format\.func\d+$`).MatchString(lines[2]) {
			t.Fatalf("first frame = %q", lines[2])
		}
	})

	t.Run("+v of a self-nested error terminates", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("loop").Yeet()
		_ = err.Nest(err)
		if got := fmt.Sprintf("%+v", err); got != "loop" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("nil and zero value", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			err  *oops.Error
			want string
		}{
			{nil, "oops.Error(nil)"},
			{&oops.Error{}, "oops.Error(undefined)"},
		} {
			for _, format := range []string{"%v", "%+v", "%s"} {
				if got := fmt.Sprintf(format, tc.err); got != tc.want {
					t.Errorf("%s of %#v: got %q, want %q", format, tc.err, got, tc.want)
				}
			}
		}
	})
}
