package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

// Thrashing optimisation
//
// Find pairs of store/load instructions (non-similar pairs are allowed) operating on the same box. If the instructions
// between them are not accumulating instructions then the second of the pair can be removed.

func thrash(prog *lmc.Program) error {
	previous := -1

	instrs := make([]lmc.Instruction, len(prog.Memory.GetInstructionSet().GetInstructions()))
	copy(instrs, prog.Memory.GetInstructionSet().GetInstructions())

	for i, removed := 0, 0; i < len(instrs); i++ {
		var ok bool

		_, ok = instrs[i].(*lmc.StoreInstr)
		if !ok {
			_, ok = instrs[i].(*lmc.LoadInstr)
		}

		if ok {
			if previous != -1 && (instrs[i].Boxes()[0].Address() != instrs[previous].Boxes()[0].Address()) {
				previous = -1
			}

			if i == len(instrs)-1 {
				if err := prog.Memory.GetInstructionSet().RemoveInstruction(i - removed); err != nil {
					return err
				} else {
					removed++
				}

				break
			}

			if previous == -1 {
				previous = i
			} else {
				remove := true

				for j := previous + 1; j < i; j++ {
					if instrs[j].ACC() {
						remove = false
						break
					}
				}

				if remove {
					if err := prog.Memory.GetInstructionSet().RemoveInstruction(i - removed); err != nil {
						return err
					} else {
						removed++
					}
				} else {
					previous = i
				}
			}
		}
	}

	return nil
}

type OThrashing struct {
	program *lmc.Program
}

func NewOThrashing(program *lmc.Program) *OThrashing {
	return &OThrashing{
		program: program,
	}
}

func (o *OThrashing) Strategy() OStrategy {
	return Thrashing
}

func (o *OThrashing) Optimise() error {
	if err := thrash(o.program); err != nil {
		return fmt.Errorf("thrashing failed: %s", err)
	}

	return nil
}

func (o *OThrashing) Program() *lmc.Program {
	return o.program
}
