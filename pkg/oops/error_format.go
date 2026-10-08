package oops

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

// Format implements fmt.Formatter. %+v prints err as a multi-line tree: the
// first line is Error(), followed by indented lines for the path, causes,
// actions, fields (sorted by key) and trace frames. Each error below it, found
// with Walk, follows on a line marked "└" and indented by its depth; a foreign
// error prints its own Error(). Every other verb formats Error() as a string,
// so %v and %s print Error() and %q prints it quoted.
func (err *Error) Format(f fmt.State, verb rune) {
	if verb != 'v' || !f.Flag('+') {
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), err.Error())
		return
	}

	if err == nil {
		_, _ = io.WriteString(f, err.Error())
		return
	}

	for depth, node := range Walk(err) {
		indent := strings.Repeat("  ", depth+1)
		header := node.Error()
		if depth > 0 {
			_, _ = io.WriteString(f, "\n"+indent[2:]+"└ ")
			header = strings.ReplaceAll(header, "\n", "\n"+indent)
		}
		_, _ = io.WriteString(f, header)

		if e, ok := node.(*Error); ok { //nolint:errorlint // each walked node is inspected directly
			e.formatAttributes(f, indent)
		}
	}
}

// formatAttributes writes err's path, causes, actions, fields and trace, one
// per line, each starting with a newline and indent.
func (err *Error) formatAttributes(w io.Writer, indent string) {
	if err.path != "" {
		_, _ = io.WriteString(w, "\n"+indent+"path: "+err.path)
	}

	if len(err.causes) > 0 {
		_, _ = io.WriteString(w, "\n"+indent+"causes: "+strings.Join(err.causes, ", "))
	}

	if len(err.actions) > 0 {
		_, _ = io.WriteString(w, "\n"+indent+"actions: "+strings.Join(err.actions, ", "))
	}

	if len(err.fields) > 0 {
		_, _ = io.WriteString(w, "\n"+indent+"fields: ")
		for idx, key := range slices.Sorted(maps.Keys(err.fields)) {
			if idx > 0 {
				_, _ = io.WriteString(w, ", ")
			}
			_, _ = fmt.Fprintf(w, "%s=%v", key, err.fields[key])
		}
	}

	if len(err.trace) > 0 {
		_, _ = io.WriteString(w, "\n"+indent+"trace:")
		for _, frame := range err.trace {
			_, _ = io.WriteString(w, "\n"+indent+"  "+frame)
		}
	}
}
