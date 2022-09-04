package optimisation

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
)

// Thrashing optimisation
//
// Find pairs of store/load instructions (non-similar pairs are allowed) operating on the same box. If the instructions
// between them are not accumulating instructions then the second of the pair can be removed.

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
	var removals []int
	previous := -1
	instrs := o.program.Memory.GetInstructionSet().GetInstructions()

	for i := 0; i < len(instrs); i++ {
		var ok bool

		_, ok = instrs[i].(*lmc.StoreInstr)
		if !ok {
			_, ok = instrs[i].(*lmc.LoadInstr)
		}

		ok = ok && (previous == -1 || (instrs[i].Boxes()[0].Address() == instrs[previous].Boxes()[0].Address()) )

		if ok {
			if i == len(instrs) - 1 {
				removals = append(removals, i)
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
					removals = append(removals, i)
				} else {
					previous = i
				}
			}
		}
	}

	for k, l := range removals {
		if err := o.program.Memory.GetInstructionSet().RemoveInstruction(l - k); err != nil {
			return err
		}
	}

	return nil
}

func (o *OThrashing) Program() *lmc.Program {
	return o.program
}
