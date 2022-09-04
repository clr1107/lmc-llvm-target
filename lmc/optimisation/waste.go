package optimisation

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
)

type OWaste struct {
	program  *lmc.Program
}

func NewOWaste(program *lmc.Program) *OWaste {
	return &OWaste{
		program: program,
	}
}

func (O *OWaste) Strategy() OStrategy {
	return Waste
}

func (O *OWaste) Optimise() error {
	used := make(map[string]int, 0)

	for _, v := range O.program.Memory.GetMailboxes() {
		used[v.Identifier()] = 0
	}

	for _, v := range O.program.Memory.GetInstructionSet().GetInstructions() {
		for _, box := range v.Boxes() {
			used[box.Identifier()]++
		}
	}

	for id, c := range used {
		if c == 0 {
			O.program.Memory.RemoveMailboxIdentifier(id)
		}
	}

	return nil
}

func (O *OWaste) Program() *lmc.Program {
	return O.program
}
