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

var (
	MailboxAlreadyExistsAddressError = func(addr Address) error {
		return fmt.Errorf("a mailbox with address %d already exists", addr)
	}
	MailboxAlreadyExistsIdentifierError = func(identifier string) error {
		return fmt.Errorf("a mailbox with identifier `%s' already exists", identifier)
	}
	LabelAlreadyExistsError = func(identifier string) error {
		return fmt.Errorf("a label with identifier `%s' already exists", identifier)
	}
	CannotRemoveInstructionIndexError = func(a int, b int) error {
		return fmt.Errorf("cannot remove instruction index %d out of %d", a, b)
	}
	VariableDoesNotExistError = func(name string) error {
		return fmt.Errorf("variable `%s` does not exist", name)
	}
)

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
