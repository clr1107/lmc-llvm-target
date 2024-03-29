package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

var thrashStageNames = [...]string{
	"THRASH_MUL_LOAD",
	"THRASH_PAIRS",
}

func thrashErr(stage int, child error) error {
	return fmt.Errorf("thrashing failed stage %d=%s: %s", stage, thrashStageNames[stage], child)
}

func thrash_mul_load(prog *lmc.Program) error {
	// Ugly... this entire function is being redone soon. But, it works.
	instrs := prog.Memory.InstructionsList.Instructions
	previous := -1

	for i, removed := 0, 0; i < len(instrs); i++ {

		if _, ok := instrs[i-removed].(*lmc.LoadInstr); !ok {
			continue
		}

		if previous == -1 {
			previous = i
			continue
		}

		var acc bool

		for j := previous + 1; j < i; j++ {
			if instrs[j-removed].ACC() {
				acc = true
				break
			}
		}

		if !acc {
			if err := prog.Memory.InstructionsList.RemoveInstruction(previous - removed); err != nil {
				return err
			} else {
				removed++
				previous = i
			}
		}

	}

	return nil
}

func thrash_pairs(prog *lmc.Program) error {
	previous := -1

	instrs := make([]lmc.Instruction, len(prog.Memory.InstructionsList.Instructions))
	copy(instrs, prog.Memory.InstructionsList.Instructions)

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
				if err := prog.Memory.InstructionsList.RemoveInstruction(i - removed); err != nil {
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
					if err := prog.Memory.InstructionsList.RemoveInstruction(i - removed); err != nil {
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
	var err error

	if err = thrash_mul_load(o.program); err != nil {
		return thrashErr(0, err)
	} else if err = thrash_pairs(o.program); err != nil {
		return thrashErr(1, err)
	}

	return nil
}

func (o *OThrashing) Program() *lmc.Program {
	return o.program
}
