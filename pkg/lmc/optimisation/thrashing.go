package optimisation

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
)

type OThrashing struct {
	program  *lmc.Program
}

func NewOThrashing(program *lmc.Program) *OThrashing {
	return &OThrashing{
		program: program,
	}
}

func (O *OThrashing) Strategy() OStrategy {
	return Thrashing
}

func (O *OThrashing) Optimise() error {
	var i []int
	instrs := O.program.Memory.GetInstructionSet().GetInstructions()

	for ii := 1; ii < len(instrs); ii++ {
		v2, ok2 := instrs[ii].(*lmc.LoadInstr)
		if ok2 {
			v1, ok1 := instrs[ii-1].(*lmc.StoreInstr)
			if ok1 && v2.Param.Identifier() == v1.Param.Identifier() {
				i = append(i, ii)
			}
		}
	}

	for k, ii := range i {
		if err := O.program.Memory.GetInstructionSet().RemoveInstruction(ii-k); err != nil {
			return err
		}
	}

	return nil
}

func (O *OThrashing) Program() *lmc.Program {
	return O.program
}
