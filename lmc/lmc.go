// Package lmc allows interaction with Little Man Computer, programmatically.
// It is restricted to static analyses and optimisations, and does not parse
// text-form LMC. (Note: this package is simple enough that it doesn't check for
// non-trivial things such as >100 mailboxes.)
//
// # Advisory note
//
// Note on memory operations: any newly created mailboxes or labels given in a
// memory operation must be manually added to any program! During any function
// the operations are presumed to work, and be completed in the future. Use the
// provided function in *Program#AddMemoryOp.
// Consider using utility functions provided in *Program, which do this for you.
package lmc

import (
	"fmt"
)

type Address int64
type Value int

type LMCType interface {
	fmt.Stringer
	LMCString() string
}

type Addressable interface {
	Address() Address
}

type Identifiable interface {
	Identifier() string
}

// ---------- Error ----------

type Error struct {
	error
	E      error
	Header string
	Child  error
}

func (e *Error) Error() string {
	var s string

	if e.Header != "" {
		s = e.Header + " "
	}

	s += e.E.Error()

	if e.Child != nil {
		s += ": " + e.Child.Error()
	}

	return s
}
