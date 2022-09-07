package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

// Box propagation
//
// Find any boxes that merely serve a temporary purpose and remove their use, replacing them with their permanent box.
// E.g., the pair STA A; STA X

var stageNames = [...]string{
	"PROP_STA_STA",
	"PROP_LDA_STA",
}

func propErr(stage int, child error) error {
	return fmt.Errorf("box propogation failed stage %d=%s: %s", stage, stageNames[stage], child)
}

func prop_sta_sta(prog *lmc.Program) error {
	instrs := make([]lmc.Instruction, len(prog.Memory.GetInstructionSet().GetInstructions()))
	copy(instrs, prog.Memory.GetInstructionSet().GetInstructions())

	for i, removed := 1, 0; i < len(instrs); i++ {
		var ok bool

		_, ok = instrs[i-1].(*lmc.StoreInstr)
		if ok {
			if _, ok2 := instrs[i].(*lmc.StoreInstr); !ok2 {
				ok = false
				i++
			}
		}

		if ok {
			if err := prog.Memory.GetInstructionSet().RemoveInstruction(i - 1 - removed); err != nil {
				return err
			} else {
				removed++
			}
		}
	}

	return nil
}

func prop_lda_sta(prog *lmc.Program) error {
	instrs := make([]lmc.Instruction, len(prog.Memory.GetInstructionSet().GetInstructions()))
	copy(instrs, prog.Memory.GetInstructionSet().GetInstructions())

	previous := -1

	for i, removed := 0, 0; i < len(instrs); i++ {
		var ok bool

		if _, ok = instrs[i].(*lmc.LoadInstr); ok {
			previous = i
			continue
		}

		ok = previous != -1
		if ok {
			if _, ok2 := instrs[i].(*lmc.StoreInstr); !ok2 {
				ok = false
			}
		}

		if ok {
			if instrs[i].Boxes()[0].Address() != instrs[previous].Boxes()[0].Address() {
				previous = -1
				continue
			}

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
			}
		}
	}

	return nil
}

type OProp struct {
	program *lmc.Program
}

func NewOProp(program *lmc.Program) *OProp {
	return &OProp{
		program: program,
	}
}

func (o *OProp) Strategy() OStrategy {
	return BProp
}

func (o *OProp) Optimise() error {
	var err error

	if err = prop_sta_sta(o.program); err != nil {
		return propErr(0, err)
	} else if err = prop_lda_sta(o.program); err != nil {
		return propErr(1, err)
	}

	return nil
}

func (o *OProp) Program() *lmc.Program {
	return o.program
}
