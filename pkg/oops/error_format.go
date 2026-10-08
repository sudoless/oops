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
// actions, fields (sorted by key) and trace frames. Each error below it, in
// Walk order, follows on a line marked "└" and indented by its depth. A
// foreign error prints its own Error() once and the errors below it are not
// printed, since its text usually already includes them. Text that spans
// several lines is indented to match. Cycles are cut as in Walk; when the
// tree exceeds Walk's 1024-node cap, the output ends with a line saying so.
//
// A definition's Formatter controls Error() and therefore only the first line
// of each *Error in the tree; the attribute and child lines are always
// rendered by oops.
//
// Every other verb formats Error() as a string, so %v and %s print Error() and
// %q prints it quoted.
func (err *Error) Format(f fmt.State, verb rune) {
	if verb != 'v' || !f.Flag('+') {
		_, _ = fmt.Fprintf(f, fmt.FormatString(f, verb), err.Error())
		return
	}

	if err == nil {
		_, _ = io.WriteString(f, err.Error())
		return
	}

	foreignDepth := -1 // depth of the foreign node whose subtree is skipped
	truncated := walk(err, func(depth int, node error) bool {
		if foreignDepth >= 0 && depth > foreignDepth {
			return true
		}
		foreignDepth = -1

		indent := strings.Repeat("  ", depth+1)
		if depth > 0 {
			_, _ = io.WriteString(f, "\n"+indent[2:]+"└ ")
		}
		_, _ = io.WriteString(f, indentLines(node.Error(), indent))

		if e, ok := node.(*Error); ok { //nolint:errorlint // each walked node is inspected directly
			e.formatAttributes(f, indent)
		} else {
			foreignDepth = depth
		}

		return true
	})

	if truncated {
		_, _ = fmt.Fprintf(f, "\n  … truncated after %d nodes", walkLimit)
	}
}

// indentLines indents every line of s after the first by indent.
func indentLines(s, indent string) string {
	return strings.ReplaceAll(s, "\n", "\n"+indent)
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
			_, _ = io.WriteString(w, key+"="+indentLines(fmt.Sprint(err.fields[key]), indent))
		}
	}

	if len(err.trace) > 0 {
		_, _ = io.WriteString(w, "\n"+indent+"trace:")
		for _, frame := range err.trace {
			_, _ = io.WriteString(w, "\n"+indent+"  "+frame)
		}
	}
}
