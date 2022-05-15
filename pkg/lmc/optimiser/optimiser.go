package optimiser

import "github.com/clr1107/lmc-llvm-target/lmc"

type OStrategy uint

const (
	Thrashing OStrategy = iota
	Waste
	BProp
	Chaining
	Unroll
)

type Optimiser interface {
	Optimise() error
	Program() *lmc.Program
	Strategy() OStrategy
}

