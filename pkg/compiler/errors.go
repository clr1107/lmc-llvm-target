package compiler

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"strings"
)

var (
	NonexistentPropertyError = func(property string) error {
		return fmt.Errorf("nonexistent propert `%s'", property)
	}
	IncorrectTypesError = func(types... string) error {
		return fmt.Errorf("incorrect types: %s", strings.Join(types, ", "))
	}
	InvalidLLTypeError = func(types... string) error {
		return fmt.Errorf("incorrect LL types: %s", strings.Join(types, ", "))
	}
	UnknownMailboxError = func(addr lmc.Address) error {
		return fmt.Errorf("unknown mailbox with address %d\n", addr)
	}
)
