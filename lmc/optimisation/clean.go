package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

var cleanStageNames = [...]string{
	"CLEAN_DEAD",
	"CLEAN_MULTI",
}

func cleanErr(stage int, child error) error {
	return fmt.Errorf("cleaning failed stage %d=%s: %s", stage, cleanStageNames[stage], child)
}

func clean_dead_box(prog *lmc.Program) error {
	used := make(map[string]int, 0)

	for _, def := range prog.Memory.GetInstructionSet().DefInstructions {
		used[def.Boxes()[0].Identifier()] = 0
	}

	for _, v := range prog.Memory.GetInstructionSet().Instructions {
		for _, box := range v.Boxes() {
			used[box.Identifier()]++
		}
	}

	for id, c := range used {
		if c == 0 {
			if err := prog.Memory.GetInstructionSet().RemoveDef(id); err != nil {
				return err
			}
			prog.Memory.RemoveMailboxIdentifier(id)
		}
	}

	return nil
}

func clean_multi_dat(prog *lmc.Program) error {
	seen := make(map[lmc.Address]struct{})
	instrs := prog.Memory.GetInstructionSet().DefInstructions

	var ok bool
	var i int

	for _, ii := range instrs {
		if _, ok = seen[ii.Box.Address()]; !ok {
			instrs[i] = ii
			i++

			seen[ii.Box.Address()] = struct{}{}
		}
	}

	for j := i; j < len(instrs); j++ {
		instrs[j] = nil
	}

	prog.Memory.GetInstructionSet().DefInstructions = instrs[:i]
	return nil
}

type OClean struct {
	program *lmc.Program
}

func NewOClean(program *lmc.Program) *OClean {
	return &OClean{
		program: program,
	}
}

func (o *OClean) Strategy() OStrategy {
	return Clean
}

func (o *OClean) Optimise() error {
	_ = clean_dead_box(o.program)  // currently no error returned
	_ = clean_multi_dat(o.program) // ^^

	return nil
}

func (o *OClean) Program() *lmc.Program {
	return o.program
}
